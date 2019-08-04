package tracking

import (
	"github.com/markcheno/go-talib"
	geo "github.com/paulmach/go.geo"
)

func getIndicatorStatuses(candles []Candle) IndicatorsStatus {
	for {
		// преобразование свечей
		closePrices := make([]float64, len(candles))
		for index, candle := range candles {
			closePrices[index] = candle.Close
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
				return IndicatorsStatusBuy
			} else if firstLineStochRsi.Intersection(secondLineStochRsi).Y() > 80 &&
				kShortSchRsi[len(kShortSchRsi)-3] > kShortSchRsi[len(kShortSchRsi)-2] &&
				dLongSchRsi[len(dLongSchRsi)-3] > dLongSchRsi[len(dLongSchRsi)-2] {
				return IndicatorsStatusSell
			}
		}

		return IndicatorsStatusNeutral
	}
}
