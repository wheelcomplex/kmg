package typeTransform

import "testing"
import (
	"github.com/bronze1man/kmg/kmgTest"
	"reflect"
)

func TestManager(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	Int := 0
	ArrMapStringInt := []map[string]int{}
	type T1 struct {
		A int
		B string
	}
	ArrStruct := []T1{}
	StringSlice := []string{}
	testCaseTable := []struct {
		in  interface{}
		out interface{}
		exp interface{}
	}{
		{1, &Int, 1}, //0
		{int64(1), &Int, 1},
		{ //2
			[]map[string]string{
				{
					"a": "1",
				},
				{
					"b": "1",
				},
			},
			&ArrMapStringInt,
			[]map[string]int{
				{
					"a": 1,
				},
				{
					"b": 1,
				},
			},
		},
		{ //3
			[]map[string]string{
				{
					"A": "1",
					"B": "abc",
					"C": "abd",
				},
				{
					"A": "",
					"B": "",
					"C": "abd",
				},
			},
			&ArrStruct,
			[]T1{
				{
					A: 1,
					B: "abc",
				},
				{
					A: 0,
					B: "",
				},
			},
		},
		{ //4
			[]interface{}{
				"1",
				"2",
			},
			&StringSlice,
			[]string{
				"1",
				"2",
			},
		},
	}
	for i, testCase := range testCaseTable {
		err := Transform(testCase.in, testCase.out)
		t.EqualMsg(err, nil, "fail at %d", i)
		t.EqualMsg(reflect.ValueOf(testCase.out).Elem().Interface(), testCase.exp, "fail at %d", i)
	}
}
