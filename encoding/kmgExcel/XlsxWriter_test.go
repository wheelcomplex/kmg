package kmgExcel

import "testing"
import "github.com/bronze1man/kmg/kmgTest"

func TestCoordinateXy2Excel(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	for _, c := range []struct {
		x   int
		y   int
		out string
	}{
		{0, 0, "A1"},
		{25, 0, "Z1"},
		{25, 1000, "Z1001"},
		{26, 1000, "AA1001"},
		{27, 1000, "AB1001"},
		{676, 0, "ZA1"},
		{701, 0, "ZZ1"},
		{702, 0, "AAA1"},
		{703, 0, "AAB1"},
		{728, 0, "ABA1"},
	} {
		t.Equal(CoordinateXy2Excel(c.x, c.y), c.out)
	}
}
