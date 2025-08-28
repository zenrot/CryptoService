package postgresStorage

import (
	"CryptoService/internal/config"
	"CryptoService/internal/crypt"
	"CryptoService/internal/storage"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"math"
	"sync"
	"time"
)

type postgresStorage struct {
	mu             sync.RWMutex
	symbToIDmap    map[string]int
	postgresConfig *config.PostgresConfig
	db             *sql.DB
}

func NewPostgresStorage(cfg *config.Config) (*postgresStorage, error) {
	pc := cfg.PostgresConfig
	var psqlInfo string
	if pc.Password == "" {
		psqlInfo = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable",
			pc.Host, pc.Port, pc.User, pc.Dbname,
		)
	} else {
		psqlInfo = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			pc.Host, pc.Port, pc.User, pc.Password, pc.Dbname,
		)
	}

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
       user_name text NOT NULL UNIQUE,
       password text NOT NULL,
	   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       PRIMARY KEY(user_name)
);`)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS crypto_info (
    crypto_id serial PRIMARY KEY,
    name text NOT NULL UNIQUE,
    symbol text NOT NULL
);`)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS crypto_prices (
    crypto_id int NOT NULL,
    price float NOT NULL,
    timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(crypto_id) REFERENCES crypto_info(crypto_id)
);`)
	if err != nil {
		return nil, err
	}
	symbToIDmap := make(map[string]int)
	rows, err := db.Query(`SELECT crypto_id, symbol FROM crypto_info WHERE crypto_id IS NOT NULL`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var symbol string
		err = rows.Scan(&id, &symbol)
		if err != nil {
			return nil, err
		}
		symbToIDmap[symbol] = id
	}
	return &postgresStorage{
		symbToIDmap:    symbToIDmap,
		postgresConfig: &cfg.PostgresConfig,
		db:             db,
	}, nil
}

func (st *postgresStorage) RegisterUser(name, password string) error {
	hashedPasswd, err := crypt.HashPassword(password)
	if err != nil {
		return err
	}
	_, err = st.db.Exec(`INSERT INTO users (user_name, password) VALUES ($1, $2)`, name, hashedPasswd)
	if err != nil {
		return storage.ErrUserExists
	}
	return nil
}

func (st *postgresStorage) LoginUser(name, password string) (*storage.User, error) {
	rows := st.db.QueryRow(`SELECT user_name, password FROM users WHERE user_name = $1`, name)
	var user storage.User
	err := rows.Scan(&user.Name, &user.Password)
	if err != nil {
		return nil, storage.ErrUserNotExists
	}
	if !crypt.CheckPasswordHash(password, user.Password) {
		return nil, storage.ErrWrongPassword
	}
	return &user, nil
}

func (st *postgresStorage) AddCrypto(symbol, name string, price float64, t time.Time) error {
	st.mu.Lock()
	defer st.mu.Unlock()
	var cryptoID int
	if val, ok := st.symbToIDmap[symbol]; !ok {
		err := st.db.QueryRow(
			`INSERT INTO crypto_info (name, symbol) VALUES ($1, $2) RETURNING crypto_id`,
			name, symbol,
		).Scan(&cryptoID)
		if err != nil {
			return err
		}
		st.symbToIDmap[symbol] = cryptoID
	} else {
		cryptoID = val
	}

	_, err := st.db.Exec(
		`INSERT INTO crypto_prices (crypto_id, price, timestamp) VALUES ($1, $2, $3)`,
		cryptoID, price, t,
	)
	if err != nil {
		return err
	}

	return nil
}

func (st *postgresStorage) GetCrypto(symbol string) ([]storage.CryptoVal, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	res := make([]storage.CryptoVal, 0)
	if id, ok := st.symbToIDmap[symbol]; !ok {
		return nil, fmt.Errorf("symbol %s is not being tracked", symbol)
	} else {
		rows, err := st.db.Query(`SELECT  cp.price, cp.timestamp, ci.name FROM crypto_prices as cp JOIN crypto_info as ci USING (crypto_id) WHERE crypto_id = $1`, id)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var value storage.CryptoVal
			if err := rows.Scan(&value.Price, &value.Time, &value.Name); err != nil {
				return nil, err
			}
			value.Symbol = symbol
			res = append(res, value)
		}
	}
	return res, nil
}

func (st *postgresStorage) GetLatestCrypto() (map[string]storage.CryptoVal, error) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	res := make(map[string]storage.CryptoVal)
	for symbol := range st.symbToIDmap {
		var value storage.CryptoVal
		value.Symbol = symbol
		err := st.db.QueryRow(`SELECT cp.price, cp.timestamp, ci.name
			FROM crypto_prices AS cp
			JOIN crypto_info AS ci USING (crypto_id)
			WHERE ci.symbol = $1
			ORDER BY cp.timestamp DESC
			LIMIT 1;`, symbol).Scan(&value.Price, &value.Time, &value.Name)
		if err != nil {
			return nil, err
		}
		res[symbol] = value
	}
	return res, nil
}

func (st *postgresStorage) DeleteCrypto(symbol string) error {
	_, err := st.db.Exec(`DELETE FROM crypto_prices cp
		USING crypto_info ci
		WHERE cp.crypto_id = ci.crypto_id
		  AND ci.symbol = $1;`, symbol)
	if err != nil {
		return err
	}
	_, err = st.db.Exec(`DELETE FROM crypto_info WHERE symbol = $1`, symbol)
	if err != nil {
		return err
	}
	delete(st.symbToIDmap, symbol)
	return nil
}

func (st *postgresStorage) GetCryptoStats(symbol string) (storage.CryptoStat, error) {
	res, err := st.GetCrypto(symbol)
	if err != nil {
		return storage.CryptoStat{}, err
	}
	if len(res) == 0 {
		return storage.CryptoStat{}, fmt.Errorf("no records for %s", symbol)
	}

	max := 0.0
	min := math.Inf(1)
	sum := 0.0

	for _, v := range res {
		if v.Price < min {
			min = v.Price
		}
		if v.Price > max {
			max = v.Price
		}
		sum += v.Price
	}

	last := res[len(res)-1].Price
	avg := sum / float64(len(res))
	priceChange := last - min
	priceChangePercent := 0.0
	if min != 0 {
		priceChangePercent = (priceChange / min) * 100
	}

	return storage.CryptoStat{
		MinPrice:           min,
		MaxPrice:           max,
		AvgPrice:           avg,
		PriceChange:        priceChange,
		PriceChangePercent: priceChangePercent,
		RecordsCount:       len(res),
	}, nil
}
