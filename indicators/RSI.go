package indicators

// RSI - индекс относительной силы
func RSI(closePrices []float64) (float64, error) {
	priceIncrease, priceReduction, size := 0.0, 0.0, len(closePrices)

	for count := 1; count < size; count++ {
		closePriceCurrent := closePrices[count]
		closePricePrev := closePrices[count-1]

		if closePriceCurrent > closePricePrev {
			priceIncrease += closePriceCurrent - closePricePrev
		} else {
			priceReduction += closePricePrev - closePriceCurrent
		}
	}

	period := float64(size)
	rs := (priceIncrease / period) / (priceReduction / period)
	rsi := 100 - 100/(1+rs)

	return rsi, nil
}
