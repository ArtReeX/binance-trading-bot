package binance

import (
	"github.com/adshao/go-binance"
)

// API - интерфейс клиента
type API struct {
	client *binance.Client
}

// NewClient - функция создания нового клиента для биржи
func NewClient(key string, secret string) *API {
	return &API{
		client: binance.NewClient(key, secret),
	}
}
