package binance

import (
	"context"
	"errors"
	"strconv"
)

// GetCurrentPrice - функция получения текущей цены валюты
func (api *API) GetCurrentPrice(pair string) (float64, error) {
	stat, err := api.client.NewPriceChangeStatsService().Symbol(pair).Do(context.Background())
	if err != nil {
		return 0, errors.New("Не удалось получить текущую цену валюты: " + err.Error())
	}
	currentPrice, err := strconv.ParseFloat(stat.LastPrice, 64)
	if err != nil {
		return 0, errors.New("Не удалось преобразовать строку текущей цены валюты в дробь: " + err.Error())
	}
	return currentPrice, nil
}
