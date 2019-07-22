package indicators

import "testing"

func TestRSI(t *testing.T) {
	closePrices := []float64{3, 2, 4, 5, 5, 7, 10, 1, 2, 3, 8, 12, 13, 10}
	rsi, _ := RSI(closePrices)
	if rsi != 60.6060606060606 {
		t.Error("Запрошено 60.6060606060606, получено", rsi)
	}
}
