package tracking

import (
	bnc "../binance"
)

func getIndicatorStatuses( /*candles []bnc.Candle,*/ depth bnc.Depth) IndicatorsStatus {

	/*
		// преобразование свечей
		closePrices := make([]float64, len(candles))
		for index, candle := range candles {
			closePrices[index] = candle.Close
		}

		// получение статуса индикаторов
		kShortSchRsi, dLongSchRsi := talib.StochRsi(closePrices, 14, 9, 3,
			3)

		firstLineStochRsi := geo.NewLine(geo.NewPoint(0, kShortSchRsi[len(kShortSchRsi)-3]),
			geo.NewPoint(1, kShortSchRsi[len(kShortSchRsi)-2]))
		secondLineStochRsi := geo.NewLine(geo.NewPoint(0, dLongSchRsi[len(dLongSchRsi)-3]),
			geo.NewPoint(1, dLongSchRsi[len(dLongSchRsi)-2]))

		if firstLineStochRsi.Intersects(secondLineStochRsi) {
			if firstLineStochRsi.Intersection(secondLineStochRsi).Y() < 20 &&
				firstLineStochRsi.Intersection(secondLineStochRsi).Y() < kShortSchRsi[len(kShortSchRsi)-2] &&
				firstLineStochRsi.Intersection(secondLineStochRsi).Y() < dLongSchRsi[len(dLongSchRsi)-2] {
				return IndicatorsStatusBuy
			} else if firstLineStochRsi.Intersection(secondLineStochRsi).Y() > 80 &&
				firstLineStochRsi.Intersection(secondLineStochRsi).Y() > kShortSchRsi[len(kShortSchRsi)-2] &&
				firstLineStochRsi.Intersection(secondLineStochRsi).Y() > dLongSchRsi[len(dLongSchRsi)-2] {
				return IndicatorsStatusSell
			}
		}

	*/

	quantityAsk, quantityBid := 0.0, 0.0
	for _, ask := range depth.Asks {
		quantityAsk += ask.Quantity
	}
	for _, bid := range depth.Bids {
		quantityBid += bid.Quantity
	}

	if quantityBid/quantityAsk > 2 {
		return IndicatorsStatusBuy
	} else if quantityAsk/quantityBid > 2 {
		return IndicatorsStatusSell
	}

	return IndicatorsStatusNeutral
}
