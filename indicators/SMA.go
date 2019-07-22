package indicators

// SMA - оссцилилятор на основе стредних значений за период
func SMA(values []float64) float64 {
	sum := 0.0
	for _, value := range values {
		sum += value
	}
	return sum / float64(len(values))
}
