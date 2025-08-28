package ramstore

import (
	"CryptoService/internal/crypt"
	"CryptoService/internal/storage"
	"CryptoService/internal/storage/ramstore/ringBuffer"
	"fmt"
	"math"
	"sync"
	"time"
)

type ramStorage struct {
	userData   map[string]storage.User
	cryptoData map[string]*ringBuffer.RingBuffer
	mu         sync.RWMutex
}

const maxHistory = 100

func NewRamStorage() *ramStorage {
	return &ramStorage{
		userData:   make(map[string]storage.User),
		cryptoData: make(map[string]*ringBuffer.RingBuffer),
	}
}

func (rs *ramStorage) RegisterUser(name, password string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if _, ok := rs.userData[name]; ok {
		return storage.ErrUserExists
	} else {
		hashedPasswd, err := crypt.HashPassword(password)
		if err != nil {
			return err
		}
		rs.userData[name] = storage.NewUser(name, hashedPasswd)
		return nil
	}
}

func (rs *ramStorage) LoginUser(name, password string) (*storage.User, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	user, ok := rs.userData[name]
	if !ok {
		return nil, storage.ErrUserNotExists
	}
	if !crypt.CheckPasswordHash(password, user.Password) {
		return nil, storage.ErrWrongPassword
	}
	return &user, nil
}

func (rs *ramStorage) AddCrypto(symbol, name string, price float64, time time.Time) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if _, ok := rs.cryptoData[symbol]; !ok {
		rs.cryptoData[symbol] = ringBuffer.NewRingBuffer(maxHistory)
	}

	val := storage.NewCryptoVal(symbol, name, price, time)

	rs.cryptoData[symbol].Add(val)
	return nil
}

func (rs *ramStorage) GetCrypto(symbol string) ([]storage.CryptoVal, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	if _, ok := rs.cryptoData[symbol]; !ok {
		return nil, fmt.Errorf("symbol %s is not being tracked", symbol)
	}

	res := rs.cryptoData[symbol].Values()
	return res, nil
}

func (rs *ramStorage) GetLatestCrypto() (map[string]storage.CryptoVal, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()
	res := make(map[string]storage.CryptoVal)
	for _, val := range rs.cryptoData {
		t, ok := val.Last()
		if !ok {
			continue
		}
		res[t.Symbol] = t
	}
	return res, nil
}

func (rs *ramStorage) DeleteCrypto(symbol string) error {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if _, ok := rs.cryptoData[symbol]; !ok {
		return storage.ErrCryptoNotExists
	}
	delete(rs.cryptoData, symbol)
	return nil
}

func (rs *ramStorage) GetCryptoStats(symbol string) (storage.CryptoStat, error) {
	rs.mu.RLock()
	defer rs.mu.RUnlock()

	res, err := rs.GetCrypto(symbol)
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
