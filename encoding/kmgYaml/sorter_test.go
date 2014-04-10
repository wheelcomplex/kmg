package kmgYaml

import "testing"
import (
	"github.com/bronze1man/kmg/kmgTest"
	"reflect"
	"sort"
)

func TestKeyListLess(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	keyList := keyList{
		reflect.ValueOf("Quest2Name_12"),  //3
		reflect.ValueOf("Quest2Name_101"), //5
	}

	testCaseTable := []struct {
		i   int
		j   int
		ret bool
	}{
		{0, 1, false},
		{1, 0, true},
	}
	for _, testCase := range testCaseTable {
		t.EqualMsg(keyList.Less(testCase.i, testCase.j), testCase.ret,
			"%d %d %v", testCase.i, testCase.j, testCase.ret)
	}
	sort.Sort(keyList)
}
