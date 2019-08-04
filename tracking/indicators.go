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
		kShortSchRsi, dLongSchRsi := talib.StochRsi(closePrices, 14, 9, 3, 3)

		firstLineStochRsi := geo.NewLine(geo.NewPoint(0, kShortSchRsi[len(kShortSchRsi)-3]),
			geo.NewPoint(1, kShortSchRsi[len(kShortSchRsi)-2]))
		secondLineStochRsi := geo.NewLine(geo.NewPoint(0, dLongSchRsi[len(dLongSchRsi)-3]),
			geo.NewPoint(1, dLongSchRsi[len(dLongSchRsi)-2]))

		if firstLineStochRsi.Intersects(secondLineStochRsi) {
			if firstLineStochRsi.Intersection(secondLineStochRsi).Y() < 20 &&
				kShortSchRsi[len(kShortSchRsi)-3] < kShortSchRsi[len(kShortSchRsi)-2] &&
				dLongSchRsi[len(dLongSchRsi)-3] < dLongSchRsi[len(dLongSchRsi)-2] {
				action <- IndicatorsStatusBuy
			} else if firstLineStochRsi.Intersection(secondLineStochRsi).Y() > 80 &&
				kShortSchRsi[len(kShortSchRsi)-3] > kShortSchRsi[len(kShortSchRsi)-2] &&
				dLongSchRsi[len(dLongSchRsi)-3] > dLongSchRsi[len(dLongSchRsi)-2] {
				action <- IndicatorsStatusSell
			}
		}

		action <- IndicatorsStatusNeutral
	}
}
