package tracking

import (
	bnc "../binance"
	"github.com/markcheno/go-talib"
	"log"
)

// trackIndicators - мониторинг индикаторов
func trackIndicators(pair string, interval Interval, client *bnc.API, action chan<- IndicatorsStatus) {
	for {
		// получение истории свечей по направлению
		candleHistory, err := client.GetCandleHistory(pair, string(interval))
		if err != nil {
			log.Println(err)
			continue
		}

		// преобразование данных
		closePrices, err := client.ConvertCandleHistory(candleHistory, bnc.Close)
		if err != nil {
			log.Println(err)
			continue
		}
		highPrices, err := client.ConvertCandleHistory(candleHistory, bnc.High)
		if err != nil {
			log.Println(err)
			continue
		}
		lowPrices, err := client.ConvertCandleHistory(candleHistory, bnc.Low)
		if err != nil {
			log.Println(err)
			continue
		}

		// получение статуса индикаторов
		rsi := talib.Rsi(closePrices, 14)
		cci := talib.Cci(highPrices, lowPrices, closePrices, 14)
		williamsR := talib.WillR(highPrices, lowPrices, closePrices, 14)
		_, _, histMACD := talib.Macd(closePrices, 12, 26, 9)

		// покупка: RSI > 50; CCI > 0; Williams %R > -20; histogramMACD > 0
		// продажа: RSI < 50; CCI < 0; Williams %R < -80; histogramMACD < 0
		if rsi[len(rsi)-2] > 50 &&
			cci[len(cci)-2] > 0 &&
			williamsR[len(williamsR)-2] > -20 &&
			histMACD[len(histMACD)-2] > 0 &&
			(rsi[len(rsi)-3] <= 50 ||
				cci[len(cci)-3] <= 0 ||
				williamsR[len(williamsR)-3] <= -20 ||
				histMACD[len(histMACD)-3] <= 0) {
			action <- IndicatorsStatusBuy
		} else if rsi[len(rsi)-2] < 50 &&
			cci[len(cci)-2] < 0 &&
			williamsR[len(williamsR)-2] < -80 &&
			histMACD[len(histMACD)-2] < 0 &&
			(rsi[len(rsi)-3] >= 50 ||
				cci[len(cci)-3] >= 0 ||
				williamsR[len(williamsR)-3] >= -80 ||
				histMACD[len(histMACD)-3] >= 0) {
			action <- IndicatorsStatusSell
		}

		action <- IndicatorsStatusNeutral
	}
}
