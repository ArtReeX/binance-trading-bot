package tracking

import (
	"github.com/adshao/go-binance"
	"strconv"
)

func FormatBinanceCandles(unformatted []*binance.Kline) []Candle {
	formatted := make([]Candle, len(unformatted))
	for index, candle := range unformatted {
		openValue, _ := strconv.ParseFloat(candle.Open, 64)
		highValue, _ := strconv.ParseFloat(candle.High, 64)
		lowValue, _ := strconv.ParseFloat(candle.Low, 64)
		closeValue, _ := strconv.ParseFloat(candle.Close, 64)
		volumeValue, _ := strconv.ParseFloat(candle.Volume, 64)
		quoteAssetVolumeValue, _ := strconv.ParseFloat(candle.QuoteAssetVolume, 64)
		takerBuyBaseAssetVolumeValue, _ := strconv.ParseFloat(candle.TakerBuyBaseAssetVolume, 64)
		takerBuyQuoteAssetVolumeValue, _ := strconv.ParseFloat(candle.TakerBuyQuoteAssetVolume, 64)

		formatted[index] = Candle{
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
	return formatted
}
