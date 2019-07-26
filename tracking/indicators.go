package tracking

import (
	"log"
	"time"

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
		kCandlePrev := kFast[len(kFast)-2]
		dCandleCurrent := dLong[len(dLong)-1]
		dCandlePrev := dLong[len(dLong)-2]

		if kCandleCurrent < 20 && dCandleCurrent < 20 {
			// если произошло пересечение быстрой прямой долгую снизу вверх в зоне перепроданности - покупка
			if kCandlePrev <= dCandlePrev && kCandleCurrent > dCandleCurrent {
				*status = IndicatorStatusBuy
			}
		} else if kCandleCurrent > 80 && dCandleCurrent > 80 {
			// если произошло пересечение быстрой прямой долгую сверху вниз в зоне перекупленности  - продажа
			if kCandlePrev >= dCandlePrev && kCandleCurrent < dCandleCurrent {
				*status = IndicatorStatusSell
			}
		} else {
			// если мы в нейтральной зоне - нейтрально
			*status = IndicatorStatusNeutral

			// если мы в нейтральной зоне увеличиваем частоту проверок
			time.Sleep(time.Second * 5)
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

		if signal[len(signal)-1] > 0 {
			// если сигнал выше нуля - покупаем
			*status = IndicatorStatusBuy
		} else if signal[len(signal)-1] < 0 {
			// если сигнал нижу нуля - продаём
			*status = IndicatorStatusSell
		} else {
			*status = IndicatorStatusNeutral
			// если сигнал на нуле - нейтрально
			time.Sleep(time.Second * 5)
		}
	}
}
