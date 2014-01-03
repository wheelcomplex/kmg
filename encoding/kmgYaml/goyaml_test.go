package kmgYaml

import "testing"
import "github.com/bronze1man/kmg/test"

func TestUnicodeMarshal(ot *testing.T) {
	t := test.NewTestTools(ot)
	testCaseTable := []struct {
		in  string
		out string
	}{
		{
			`中文`, "中文\n",
		},
	}
	for _, testCase := range testCaseTable {
		outByte, err := Marshal(testCase.in)
		t.Equal(err, nil)
		t.Equal(string(outByte), testCase.out)
	}
}
