package indicators

import (
	"errors"
)

// minRSI - функция получения минимального RSI за определённый период
func minRSI(closePrices []float64, rsiPeriod int) (float64, error) {
	minRSI, size := 9999999999999.9, len(closePrices)

	if rsiPeriod > size {
		return 0, errors.New("Размер истории валюты меньше необходимого для вычисления минимального RSI за период")
	}

	for count := 0; count < size && count+rsiPeriod+1 <= size; count++ {
		rsi, err := RSI(closePrices[count : count+rsiPeriod+1])
		if err != nil {
			return 0, err
		}

		if rsi < minRSI {
			minRSI = rsi
		}
	}
	return minRSI, nil
}

// maxRSI - функция получения максимального RSI за определённый период
func maxRSI(closePrices []float64, rsiPeriod int) (float64, error) {
	maxRSI, size := 0.0, len(closePrices)

	if rsiPeriod > size {
		return 0, errors.New("Размер истории валюты меньше необходимого для вычисления максимального RSI за период")
	}

	for count := 0; count < size && count+rsiPeriod+1 <= size; count++ {
		rsi, err := RSI(closePrices[count : count+rsiPeriod+1])
		if err != nil {
			return 0, err
		}

		if rsi > maxRSI {
			maxRSI = rsi
		}
	}
	return maxRSI, nil
}

// StochRSI - стохастический осциллилятор
func StochRSI(closePrices []float64, rsiPeriod int, stochPeriod int, kSmothing int, dSmothing int) ([]float64, []float64, error) {
	if rsiPeriod+stochPeriod+kSmothing > len(closePrices) {
		return nil, nil, errors.New("Размер истории валюты меньше необходимого для вычисления StochRSI")
	}

	// получение не сглаженных значений
	kNotSmoothed := make([]float64, len(closePrices)-stochPeriod-rsiPeriod)
	for count := 0; count < len(closePrices) && count+stochPeriod+rsiPeriod+1 <= len(closePrices); count++ {
		currentRSI, err := RSI(closePrices[count+stochPeriod : count+stochPeriod+rsiPeriod+1])
		if err != nil {
			return nil, nil, errors.New("Невозможно получить текущий RSI: " + err.Error())
		}
		minRSI, err := minRSI(closePrices[count:count+stochPeriod+rsiPeriod+1], rsiPeriod)
		if err != nil {
			return nil, nil, errors.New("Невозможно получить минимальный RSI за период: " + err.Error())
		}
		maxRSI, err := maxRSI(closePrices[count:count+stochPeriod+rsiPeriod+1], rsiPeriod)
		if err != nil {
			return nil, nil, errors.New("Невозможно получить максимальный RSI за период: " + err.Error())
		}
		kNotSmoothed[count] = (currentRSI - minRSI) / (maxRSI - minRSI) * 100
	}

	// сглаживание быстрых значений через SMA
	kSmoothed := make([]float64, len(kNotSmoothed))
	copy(kSmoothed, kNotSmoothed[0:kSmothing-1])

	for count := 0; count+kSmothing <= len(kNotSmoothed); count++ {
		kSmoothed[count+kSmothing-1] = SMA(kNotSmoothed[count : count+kSmothing])
	}

	// сглаживание медленных значений через SMA
	dSmoothed := make([]float64, len(kSmoothed))
	copy(dSmoothed, kSmoothed[0:dSmothing-1])

	for count := 0; count+dSmothing <= len(kSmoothed); count++ {
		dSmoothed[count+dSmothing-1] = SMA(kSmoothed[count : count+dSmothing])
	}

	return kSmoothed, dSmoothed, nil
}
