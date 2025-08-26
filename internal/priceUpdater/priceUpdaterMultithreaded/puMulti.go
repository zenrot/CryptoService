package priceUpdaterMultithreaded

import (
	"CryptoService/internal/config"
	"CryptoService/internal/priceUpdater"
	"CryptoService/internal/storage"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type priceUpdaterInternal struct {
	apiKey     string
	numWorkers int
	autoUpdate time.Duration
	store      storage.Storage
	coins      map[string]priceUpdater.CoinInfo
	lastUpdate time.Time

	chErrorWorkers        chan error
	chErrorSearcherDaemon chan error
	chUpdate              map[string]chan time.Duration
	chCoins               chan string
	chDelete              map[string]chan struct{}
	chRefresh             map[string]chan struct{}
	wg                    sync.WaitGroup
	mu                    sync.RWMutex
}

func NewPriceUpdaterMultithreaded(cfg *config.Config, store storage.Storage) *priceUpdaterInternal {
	return &priceUpdaterInternal{
		apiKey:     cfg.CoingeckoKey,
		autoUpdate: 3 * time.Second,
		store:      store,
	}
}

const addr = "https://api.coingecko.com"

func (pu *priceUpdaterInternal) Start() {
	pu.chErrorWorkers = make(chan error)
	pu.chErrorSearcherDaemon = make(chan error)
	pu.chDelete = make(map[string]chan struct{})
	pu.chRefresh = make(map[string]chan struct{})
	pu.chUpdate = make(map[string]chan time.Duration)
	pu.chCoins = make(chan string)
	pu.coins = make(map[string]priceUpdater.CoinInfo)
	pu.numWorkers = 0
	go pu.searcherDaemon()
	go pu.errorHandler()
}

func (pu *priceUpdaterInternal) RefreshPrice(Symbol string) error {

	if _, ok := pu.coins[Symbol]; !ok {
		return fmt.Errorf("symbol %s is not being tracked", Symbol)
	}
	pu.wg.Add(1)
	pu.chRefresh[Symbol] <- struct{}{}
	pu.wg.Wait()
	return nil
}
func (pu *priceUpdaterInternal) RefreshAllPrices() (int, error) {
	for _, coin := range pu.coins {
		pu.wg.Add(1)
		pu.chRefresh[coin.Symbol] <- struct{}{}
	}
	pu.wg.Wait()
	return pu.numWorkers, nil
}

func (pu *priceUpdaterInternal) AddCryptoTracking(Symbol string) error {
	pu.wg.Add(1)
	pu.chCoins <- Symbol
	pu.wg.Wait()
	return <-pu.chErrorSearcherDaemon
}

func (pu *priceUpdaterInternal) DeleteCryptoTracking(Symbol string) error {
	if _, ok := pu.coins[Symbol]; !ok {
		return fmt.Errorf("symbol %s is not being tracked", Symbol)
	}
	pu.chDelete[Symbol] <- struct{}{}
	delete(pu.coins, Symbol)
	return nil
}

func (pu *priceUpdaterInternal) StopUpdating() error {
	pu.wg.Add(pu.numWorkers)
	for _, coin := range pu.coins {
		pu.chUpdate[coin.Symbol] <- time.Duration(0)
	}
	pu.wg.Wait()
	if val := pu.GetUpdateTime(); val != 0*time.Second {
		return fmt.Errorf("updating was not stopped")
	}
	return nil
}

func (pu *priceUpdaterInternal) GetUpdateTime() time.Duration {
	return pu.autoUpdate
}

func (pu *priceUpdaterInternal) GetLastUpdated() time.Time {
	return pu.lastUpdate
}

func (pu *priceUpdaterInternal) ChangeUpdateTime(t time.Duration) error {
	pu.wg.Add(pu.numWorkers)
	for _, coin := range pu.coins {
		pu.chUpdate[coin.Symbol] <- t
	}
	pu.wg.Wait()
	if val := pu.GetUpdateTime(); val != t*time.Second {
		return fmt.Errorf("update time was not changed")
	}
	return nil
}

func (pu *priceUpdaterInternal) work(coin priceUpdater.CoinInfo) {
	pu.getPrice(coin)
	pu.wg.Done()
	ticker := time.NewTicker(pu.autoUpdate)
	defer ticker.Stop()
	for {
		select {
		case <-getTickerChan(ticker):
			pu.getPrice(coin)
		case t := <-pu.chUpdate[coin.Symbol]:
			if ticker != nil {
				ticker.Stop()
			}
			if t > 0 {
				ticker = time.NewTicker(t * time.Second)
				pu.autoUpdate = t * time.Second
			} else {
				ticker = nil
				pu.autoUpdate = t
			}
			pu.wg.Done()
		case <-pu.chRefresh[coin.Symbol]:
			pu.getPrice(coin)
			pu.wg.Done()
		case <-pu.chDelete[coin.Symbol]:
			close(pu.chDelete[coin.Symbol])
			close(pu.chUpdate[coin.Symbol])
			close(pu.chRefresh[coin.Symbol])
			pu.numWorkers--
			return
		}
	}
}
func getTickerChan(t *time.Ticker) <-chan time.Time {
	if t == nil {
		return nil
	}
	return t.C
}
func (pu *priceUpdaterInternal) searcherDaemon() {
	for val := range pu.chCoins {

		if _, ok := pu.coins[val]; ok {
			pu.wg.Done()
			pu.chErrorSearcherDaemon <- fmt.Errorf("this coin already exists: %s", val)
			continue
		}

		var pathInfo = fmt.Sprintf("/api/v3/search?query=%s", strings.ToLower(val))

		resp, err := http.Get(addr + pathInfo)
		if err != nil {
			pu.wg.Done()
			pu.chErrorSearcherDaemon <- err
			continue
		}

		type searchResponse struct {
			Coins []priceUpdater.CoinInfo `json:"coins"`
		}

		var searchRes searchResponse
		if err := json.NewDecoder(resp.Body).Decode(&searchRes); err != nil {
			pu.wg.Done()
			pu.chErrorSearcherDaemon <- err
			continue
		}
		resp.Body.Close()
		fl := false
		for _, coin := range searchRes.Coins {
			if coin.Symbol == val {
				fl = true
				pu.numWorkers++
				pu.chUpdate[coin.Symbol] = make(chan time.Duration)
				pu.chRefresh[coin.Symbol] = make(chan struct{})
				pu.chDelete[coin.Symbol] = make(chan struct{})
				pu.coins[val] = coin
				go pu.work(coin)
			}
		}
		if fl == false {
			pu.wg.Done()
			pu.chErrorSearcherDaemon <- fmt.Errorf("there is no coin: %s", val)
			continue
		}
		pu.chErrorSearcherDaemon <- nil
	}
}

func (pu *priceUpdaterInternal) getPrice(coin priceUpdater.CoinInfo) {
	var pathPrice = fmt.Sprintf("/api/v3/simple/price?ids=%s&vs_currencies=usd&x_cg_demo_api_key=%s",
		coin.ID, pu.apiKey)
	resp, err := http.Get(addr + pathPrice)
	if err != nil {
		pu.chErrorWorkers <- fmt.Errorf("worker %s: %q", coin.Symbol, err)
		return
	}

	var prices map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&prices); err != nil {
		pu.chErrorWorkers <- fmt.Errorf("worker %s: %q", coin.Symbol, err)
		return
	}
	price := prices[coin.ID]["usd"]
	if err := pu.store.AddCrypto(coin.Symbol, coin.Name, price, time.Now()); err != nil {
		pu.chErrorWorkers <- fmt.Errorf("worker %s: %q", coin.Symbol, err)
		return
	}
	resp.Body.Close()
	pu.mu.Lock()
	pu.lastUpdate = time.Now()
	pu.mu.Unlock()
}

func (pu *priceUpdaterInternal) errorHandler() {
	for val := range pu.chErrorWorkers {
		fmt.Println(val)
	}
}
