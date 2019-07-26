package tracking

import (
	"log"

	"github.com/markcheno/go-talib"

	bnc "../binance"
)

// trackStochRSI - мониторинг индикатора StochRSI
func trackStochRSI(pair string, interval string, status *IndicatorStatus, client *bnc.API) {
	for {
		// получение истории торгов по валюте
		candleHistory, err := client.GetCandleHistory(pair, interval)
		if err != nil {
			log.Println(err)
			continue
		}

		// преобразование данных для StochRSI
		closePrices, err := client.ConvertCandleHistory(candleHistory, bnc.Close)
		if err != nil {
			log.Println(err)
			continue
		}

		// получение данных индикатора StochRSI
		kFast, dLong := talib.StochRsi(closePrices, 14, 9, 3, talib.EMA)
		if err != nil {
			log.Println(err)
			continue
		}

		// последние две свечи
		kCandleCurrent := kFast[len(kFast)-1]
		dCandleCurrent := dLong[len(dLong)-1]

		if kCandleCurrent < 20 && dCandleCurrent < 20 {
			// если обе линии зоне перекупленности - покупка
			*status = IndicatorStatusBuy
		} else if kCandleCurrent > 80 && dCandleCurrent > 80 {
			// если обе линии зоне перекупленности - продажа
			*status = IndicatorStatusSell
		} else {
			// если мы в нейтральной зоне - нейтрально
			*status = IndicatorStatusNeutral
		}
	}
}

// trackMACD - мониторинг индикатора MACD
func trackMACD(pair string, interval string, status *IndicatorStatus, client *bnc.API) {
	for {
		// получение истории торгов по валюте
		candleHistory, err := client.GetCandleHistory(pair, interval)
		if err != nil {
			log.Println(err)
			continue
		}

		// преобразование данных для MACD
		closePrices, err := client.ConvertCandleHistory(candleHistory, bnc.Close)
		if err != nil {
			log.Println(err)
			continue
		}

		// получение данных индикатора MACD
		_, signal, _ := talib.Macd(closePrices, 12, 26, 9)
		if err != nil {
			log.Println(err)
			continue
		}

		currentSignal := signal[len(signal)-1]

		if currentSignal > 0 {
			// если сигнал выше нуля - покупаем
			*status = IndicatorStatusBuy
		} else if currentSignal < 0 {
			// если сигнал нижу нуля - продаём
			*status = IndicatorStatusSell
		} else {
			// если сигнал на нуле - нейтрально
			*status = IndicatorStatusNeutral
		}
	}
}
