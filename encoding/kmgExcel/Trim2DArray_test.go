package kmgExcel

import "testing"
import "github.com/bronze1man/kmg/kmgTest"

func TestManager(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	testCaseTable := []struct {
		in  [][]string
		out [][]string
	}{
		{
			[][]string{{}}, [][]string{},
		},
		{
			[][]string{{" ", " "}, {""}}, [][]string{},
		},
		{
			[][]string{
				{"1", "2", " "},
				{"1", " ", " "},
				{" ", " ", " "},
				{"1"},
				{" ", " "},
			},
			[][]string{
				{"1", "2"},
				{"1", " "},
				{"1", ""},
			},
		},
	}
	for _, testCase := range testCaseTable {
		ret := Trim2DArray(testCase.in)
		t.Equal(ret, testCase.out)
	}
}
