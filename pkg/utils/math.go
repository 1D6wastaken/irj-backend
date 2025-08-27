package utils

import (
	"math"

	"golang.org/x/exp/constraints"
)

type Number interface {
	constraints.Integer | constraints.Float
}

func Avg[N Number](n ...N) N {
	if len(n) == 0 {
		return N(0)
	}

	var sum N

	for _, v := range n {
		sum += v
	}

	return sum / N(len(n))
}

func RoundFloat64(value float64, precision int64) float64 {
	//nolint:gomnd
	ratio := math.Pow(10, float64(precision))

	return math.Round(value*ratio) / ratio
}
