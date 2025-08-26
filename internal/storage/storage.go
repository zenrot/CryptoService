package storage

import (
	"errors"
	"time"
)

type User struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

type CryptoVal struct {
	Symbol string    `json:"symbol"`
	Name   string    `json:"name"`
	Price  float64   `json:"price"`
	Time   time.Time `json:"time"`
}

type CryptoStat struct {
	MinPrice           float64 `json:"min_price"`
	MaxPrice           float64 `json:"max_price"`
	AvgPrice           float64 `json:"avg_price"`
	PriceChange        float64 `json:"price_change"`
	PriceChangePercent float64 `json:"price_change_percent"`
	RecordsCount       int     `json:"records_count"`
}

type Storage interface {
	RegisterUser(name, password string) error
	LoginUser(name, password string) (*User, error)
	AddCrypto(symbol, name string, price float64, time time.Time) error
	GetCrypto(symbol string) ([]CryptoVal, error)
	DeleteCrypto(symbol string) error
	GetLatestCrypto() (map[string]CryptoVal, error)
	GetCryptoStats(symbol string) (CryptoStat, error)
}

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotExists   = errors.New("user does not exists")
	ErrWrongPassword   = errors.New("wrong password")
	ErrCryptoExists    = errors.New("crypto already exists")
	ErrCryptoNotExists = errors.New("crypto does not exists")
)

func NewUser(name, password string) User {
	return User{
		Name:     name,
		Password: password,
	}
}

func NewCryptoVal(symbol, name string, price float64, time time.Time) CryptoVal {
	return CryptoVal{
		Symbol: symbol,
		Name:   name,
		Price:  price,
		Time:   time,
	}
}
