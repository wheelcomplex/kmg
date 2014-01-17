package kmgMath

import "testing"
import "github.com/bronze1man/kmg/test"

func TestFloatCompare(ot *testing.T) {
	t := test.NewTestTools(ot)
	t.Equal(Float64Equal(1.0, 1.0), true)
	t.Equal(Float64Equal(1.0/3.0*3.0, 1.0), true)
}
