package kmgReflect

import (
	"github.com/bronze1man/kmg/test"
	"reflect"
	"testing"
)

type T1 struct {
	T3
	*T4
	B int
}

type T2 struct {
	A int
	B int
	C int
}

type T3 struct {
	A int
	B int
	T2
}
type T4 struct {
	A int
	D int
}

type T5 struct {
	T6
	A int
}
type T6 int

func TestStructGetAllFieldMap(ot *testing.T) {
	t := test.NewTestTools(ot)
	t1 := reflect.TypeOf(&T1{})
	ret := StructGetAllFieldMap(t1)
	t.Equal(ret["A"].Index, []int{0, 0})
	t.Equal(ret["B"].Index, []int{2})
	t.Equal(ret["C"].Index, []int{0, 2, 2})
	t.Equal(ret["D"].Index, []int{1, 1})
	t.Equal(len(ret), 7)

	ret = StructGetAllFieldMap(reflect.TypeOf(&T5{}))
	t.Equal(ret["A"].Index, []int{1})
	t.Equal(len(ret), 2)
}
