package binance

import (
	"errors"
	"github.com/adshao/go-binance"
)

func NewClient(key string, secret string) (*Api, error) {
	client := binance.NewClient(key, secret)

	pairs, err := getPairs(client)
	if err != nil {
		return nil, errors.New("невозможно получить параметры пар для ордеров:" + err.Error())
	}

	return &Api{
		Client: client,
		Pairs:  pairs,
	}, nil
}
