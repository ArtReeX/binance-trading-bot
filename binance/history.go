package binance

import (
	"context"
	"errors"
	"github.com/adshao/go-binance"
)

// GetCandleHistory - функция получения истории цены для валюты
func (api *API) GetCandleHistory(pair string, interval string) ([]*binance.Kline, error) {
	priceHistory, err := api.client.NewKlinesService().Symbol(pair).Interval(interval).Do(context.Background())
	if err != nil {
		return nil, errors.New("не удалось получить историю валюты: " + err.Error())
	}
	return priceHistory, nil
}
