package binance

import (
	"context"
	"errors"
	"strconv"
)

// GetCandleHistory - функция получения истории цены для валюты
func (api *API) GetCandleHistory(pair string, interval string) ([]Candle, error) {
	priceHistory, err := api.Client.NewKlinesService().Symbol(pair).Interval(interval).Do(context.Background())
	if err != nil {
		return nil, errors.New("не удалось получить историю валюты: " + err.Error())
	}

	priceHistoryFormatted := make([]Candle, len(priceHistory))
	for index, candle := range priceHistory {
		openValue, _ := strconv.ParseFloat(candle.Open, 64)
		highValue, _ := strconv.ParseFloat(candle.High, 64)
		lowValue, _ := strconv.ParseFloat(candle.Low, 64)
		closeValue, _ := strconv.ParseFloat(candle.Close, 64)
		volumeValue, _ := strconv.ParseFloat(candle.Volume, 64)
		quoteAssetVolumeValue, _ := strconv.ParseFloat(candle.QuoteAssetVolume, 64)
		takerBuyBaseAssetVolumeValue, _ := strconv.ParseFloat(candle.TakerBuyBaseAssetVolume, 64)
		takerBuyQuoteAssetVolumeValue, _ := strconv.ParseFloat(candle.TakerBuyQuoteAssetVolume, 64)

		priceHistoryFormatted[index] = Candle{
			Open:                     openValue,
			High:                     highValue,
			Low:                      lowValue,
			Close:                    closeValue,
			Volume:                   volumeValue,
			QuoteAssetVolume:         quoteAssetVolumeValue,
			TakerBuyBaseAssetVolume:  takerBuyBaseAssetVolumeValue,
			TakerBuyQuoteAssetVolume: takerBuyQuoteAssetVolumeValue,
		}
	}

	return priceHistoryFormatted, nil
}
