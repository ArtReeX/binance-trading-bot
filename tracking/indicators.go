package tracking

import (
	bnc "../binance"
	"github.com/markcheno/go-talib"
	geo "github.com/paulmach/go.geo"
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

		// получение статуса индикаторов
		kShortStochRsi, dLongStochRsi := talib.StochRsi(closePrices, 14, 9, 3, 3)

		firstLineStochRsi := geo.NewLine(geo.NewPoint(0, kShortStochRsi[len(kShortStochRsi)-2]),
			geo.NewPoint(1, kShortStochRsi[len(kShortStochRsi)-1]))
		secondLineStochRsi := geo.NewLine(geo.NewPoint(0, dLongStochRsi[len(dLongStochRsi)-2]),
			geo.NewPoint(1, dLongStochRsi[len(dLongStochRsi)-1]))

		if firstLineStochRsi.Intersects(secondLineStochRsi) {
			if firstLineStochRsi.Intersection(secondLineStochRsi).Y() < 20 {
				action <- IndicatorsStatusBuy
			} else if firstLineStochRsi.Intersection(secondLineStochRsi).Y() > 80 {
				action <- IndicatorsStatusSell
			}
		}

		action <- IndicatorsStatusNeutral
	}
}
