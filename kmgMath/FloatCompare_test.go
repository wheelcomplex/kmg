package kmgMath

import "testing"
import "github.com/bronze1man/kmg/test"

func TestFloatCompare(ot *testing.T) {
	t := test.NewTestTools(ot)
	for i, testCase := range []struct {
		f      func(x float64, y float64) bool
		x      float64
		y      float64
		result bool
	}{
		{Float64LessThan, 1.0, 2.0, true},
		{Float64LessThan, 2.0, 1.0, false},
		{Float64LessThan, 2.0, 2.0, false},
		{Float64LessEqualThan, 1.0, 2.0, true},
		{Float64LessEqualThan, 1.0, 1.0, true},
		{Float64LessEqualThan, 2.0, 1.0, false}, //5

		{Float64GreaterThan, 1.0, 2.0, false},
		{Float64GreaterThan, 2.0, 1.0, true},
		{Float64GreaterThan, 1.0, 1.0, false},
		{Float64Equal, 1.0, 1.0, true},
		{Float64Equal, 1.0, 2.0, false}, //10

		{Float64Equal, 1.0 / 3.0 * 3.0, 1.0, true},
		{Float64GreaterEqualThan, 1.0, 2.0, false},
		{Float64GreaterEqualThan, 2.0, 1.0, true},
		{Float64GreaterEqualThan, 1.0, 1.0, true},
	} {
		t.EqualMsg(testCase.f(testCase.x, testCase.y), testCase.result,
			"fail at %d", i)
	}
}
