package binance

import (
	"errors"

	"github.com/adshao/go-binance"
)

// NewClient - создание нового клиента
func NewClient(key string, secret string) (*API, error) {
	client := binance.NewClient(key, secret)

	pairs, err := GetPairs(client)
	if err != nil {
		return nil, errors.New("невозможно получить параметры пар для ордеров:" + err.Error())
	}

	return &API{
		Client: client,
		Pairs:  pairs,
	}, nil
}
