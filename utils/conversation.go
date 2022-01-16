package utils

import (
	"fmt"
	"strconv"
)

func Round(x float64) float64 {
	//unit := 0.05
	//return math.Round(x/unit) * unit

	str := fmt.Sprintf("%.2f", x)
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		return f
	}

	return 0.0
}
