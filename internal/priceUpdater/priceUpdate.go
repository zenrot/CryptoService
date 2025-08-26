package priceUpdater

import (
	"time"
)

type CoinInfo struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type PriceUpdater interface {
	Start()
	RefreshPrice(Symbol string) error
	AddCryptoTracking(Symbol string) error
	DeleteCryptoTracking(Symbol string) error
	GetUpdateTime() time.Duration
	ChangeUpdateTime(t time.Duration) error
	StopUpdating() error
	GetLastUpdated() time.Time
	RefreshAllPrices() (int, error)
}
