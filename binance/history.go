package binance

import (
	"context"
	"errors"
	"strconv"

	"github.com/adshao/go-binance"
)

// TypeOfReceiving - тип получаемого значения от преобразования истории валюты
type TypeOfReceiving uint

const (
	// Close - цена закрытия свечи
	Close TypeOfReceiving = iota
	// High -  максимальная цена свечи
	High TypeOfReceiving = iota
	// Low - минимальная цена свеча
	Low TypeOfReceiving = iota
)

// GetCandleHistory - функция получения истории цены для валюты
func (api *API) GetCandleHistory(symbol string, interval string) ([]*binance.Kline, error) {
	priceHistory, err := api.client.NewKlinesService().Symbol(symbol).Interval(interval).Do(context.Background())
	if err != nil {
		return nil, errors.New("не удалось получить историю валюты: " + err.Error())
	}

	return priceHistory, nil
}

// ConvertCandleHistory - функция преобразования истории валюты
func (api *API) ConvertCandleHistory(history []*binance.Kline, typeOfReceiving TypeOfReceiving) ([]float64, error) {
	switch typeOfReceiving {
	case High:
		{
			highPrices := make([]float64, len(history))
			for index, candle := range history {
				highPrice, err := strconv.ParseFloat(candle.High, 64)
				if err != nil {
					return nil, errors.New("невозможно преобразовать строку максимальной цены свечи в дробь")
				}
				highPrices[index] = highPrice
			}
			return highPrices, nil
		}
	case Low:
		{
			lowPrices := make([]float64, len(history))
			for index, candle := range history {
				lowPrice, err := strconv.ParseFloat(candle.Low, 64)
				if err != nil {
					return nil, errors.New("невозможно преобразовать строку минимальной цены свечи в дробь")
				}
				lowPrices[index] = lowPrice
			}
			return lowPrices, nil
		}
	case Close:
		{
			closePrices := make([]float64, len(history))
			for index, candle := range history {
				closePrice, err := strconv.ParseFloat(candle.Close, 64)
				if err != nil {
					return nil, errors.New("невозможно преобразовать строку цены закрытия свечи в дробь")
				}
				closePrices[index] = closePrice
			}
			return closePrices, nil
		}
	}
	return []float64{}, nil
}
