package hmm

import "math"

func logAdd(x, y float64) float64 {
	if x == math.Inf(-1) {
		return y
	}
	if y == math.Inf(-1) {
		return x
	}
	if x > y {
		return x + math.Log1p(math.Exp(y-x))
	}
	return y + math.Log1p(math.Exp(x-y))
}

func sum(values []float64) float64 {
	result := 0.0
	for _, value := range values {
		result += value
	}
	return result
}
