package utils

import (
	"math"
)

func RoundFloat64(val float64) float64 {
	return math.Round(val*100) / 100
}

func RoundFloat32(val float32) float32 {
	return float32(RoundFloat64(float64(val)))
}
