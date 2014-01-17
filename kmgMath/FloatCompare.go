package kmgMath

import (
	"math"
)

const Epsilon = 1e-10

func Float64LessThan(x float64, y float64) bool {
	return x < (y - Epsilon)
}
func Float64LessEqualThan(x float64, y float64) bool {
	return x < (y + Epsilon)
}
func Float64GreaterThan(x float64, y float64) bool {
	return x > (y + Epsilon)
}
func Float64GreaterEqualThan(x float64, y float64) bool {
	return x > (y - Epsilon)
}
func Float64Equal(x float64, y float64) bool {
	diff := x - y
	return math.Abs(diff) < Epsilon
}
