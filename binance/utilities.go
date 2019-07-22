package binance

import (
	"errors"
	"strconv"

	"github.com/adshao/go-binance"
)

// TypeOfReceiving - тип получаемого значения от преобразования истории валюты
type TypeOfReceiving uint

const (
	// Close - цена закрытия свечи
	Close TypeOfReceiving = iota
)

// ConvertCandleHistory - функция преобразования истории валюты
func (api *API) ConvertCandleHistory(history []*binance.Kline, typeOfReceiving TypeOfReceiving) ([]float64, error) {
	switch typeOfReceiving {
	case Close:
		{
			closePrices := make([]float64, len(history))
			for index, candle := range history {
				close, err := strconv.ParseFloat(candle.Close, 64)
				if err != nil {
					return nil, errors.New("невозможно преобразовать строку закрытия свечи в дробь")
				}
				closePrices[index] = close
			}
			return closePrices, nil
		}
	}
	return []float64{}, nil
}
