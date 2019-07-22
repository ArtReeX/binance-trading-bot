package indicators

import "testing"

func TestSMA(t *testing.T) {
	average := SMA([]float64{1, 3, 5})
	if average != 3 {
		t.Error("Запрошено 3, получено", average)
	}
}
