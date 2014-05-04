package LcgCalculateConstants

import (
	"testing"
	"github.com/bronze1man/kmg/kmgRand"
	"github.com/bronze1man/kmg/kmgTest"
)

func TestLcgCalculateConstants(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	r := uint64(1e6)
	a, c := LcgCalculateConstants(r, 0)
	f := kmgRand.LcgTransformer{
		Start: 0,
		Range: r,
		A:     a,
		C:     c,
	}
	hasValueMap := make([]bool, r)
	for i := uint64(0); i < r; i++ {
		out, err := f.GenerateInRange(i)
		t.Equal(err, nil)
		hasValueMap[int(out)] = true
	}
	for i := uint64(0); i < r; i++ {
		t.Equal(hasValueMap[i], true)
	}
}
