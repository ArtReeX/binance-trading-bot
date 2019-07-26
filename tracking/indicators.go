package tracking

import (
	"log"
	"time"

	"github.com/adshao/go-binance"

	bnc "../binance"
	indicators "../indicators"
)

// trackStochRSI индикатор покупки либо продажи валюты
func trackStochRSI(pair string, interval string, action chan<- binance.SideType, client *bnc.API) {
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
		k, d, err := indicators.StochRSI(closePrices, 14, 9, 3, 3)
		if err != nil {
			log.Println(err)
			continue
		}

		// последние две свечи
		kCandleCurrent := k[len(k)-1]
		kCandlePrev := k[len(k)-2]
		dCandleCurrent := d[len(d)-1]
		dCandlePrev := d[len(d)-2]

		if kCandleCurrent < 20 && dCandleCurrent < 20 {
			// если произошло пересечение быстрой прямой долгую снизу вверх в зоне перепроданности то выполняем покупку
			if kCandlePrev <= dCandlePrev && kCandleCurrent > dCandleCurrent {
				action <- binance.SideTypeBuy
			}
			// если мы в зоне перекупленности уменьшаем частоту проверок
			time.Sleep(time.Second / 2)
		} else if kCandleCurrent > 80 && dCandleCurrent > 80 {
			// если произошло пересечение быстрой прямой долгую сверху вниз в зоне перекупленности то выполняем продажу
			if kCandlePrev >= dCandlePrev && kCandleCurrent < dCandleCurrent {
				action <- binance.SideTypeSell
			}
			// если мы в зоне перепроданности уменьшаем частоту проверок
			time.Sleep(time.Second / 2)
		} else {
			// если мы в нейтральной зоне увеличиваем частоту проверок
			time.Sleep(time.Second * 5)
		}
	}
}
