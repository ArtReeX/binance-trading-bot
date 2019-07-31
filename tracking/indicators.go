package tracking

import (
	"github.com/markcheno/go-talib"
	"log"
	"time"

	bnc "../binance"
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

		// получение статусов индикаторов
		rsi := talib.Rsi(closePrices, 14)
		cci := talib.Cci(highPrices, lowPrices, closePrices, 14)
		williamsR := talib.WillR(highPrices, lowPrices, closePrices, 14)
		_, _, histMACD := talib.Macd(closePrices, 12, 26, 9)

		if rsi[len(rsi)-2] > 50 &&
			cci[len(cci)-2] > 0 &&
			williamsR[len(williamsR)-2] > -20 &&
			histMACD[len(histMACD)-2] > 0 {
			action <- IndicatorsStatusBuy
		} else if rsi[len(rsi)-2] < 50 &&
			cci[len(cci)-2] < 0 &&
			williamsR[len(williamsR)-2] < -80 &&
			histMACD[len(histMACD)-2] < 0 {
			action <- IndicatorsStatusSell
		}

		action <- IndicatorsStatusNeutral
		time.Sleep(time.Second / 5)
	}
}
