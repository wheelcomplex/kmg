package kmgReflect

import (
	"github.com/bronze1man/kmg/test"
	"reflect"
	"testing"
)

type GetAllFieldT1 struct {
	GetAllFieldT3
	*GetAllFieldT4
	B int
}

type GetAllFieldT2 struct {
	A int
	B int
	C int
}

type GetAllFieldT3 struct {
	A int
	B int
	GetAllFieldT2
}
type GetAllFieldT4 struct {
	A int
	D int
}

type GetAllFieldT5 struct {
	GetAllFieldT6
	A int
}
type GetAllFieldT6 int

func TestStructGetAllField(ot *testing.T) {
	t := test.NewTestTools(ot)
	t1 := reflect.TypeOf(&GetAllFieldT1{})
	ret := StructGetAllField(t1)
	t.Equal(len(ret), 7)
	t.Equal(ret[0].Name, "GetAllFieldT3")
	t.Equal(ret[1].Name, "GetAllFieldT4")
	t.Equal(ret[2].Name, "B")
	t.Equal(ret[2].Index, []int{2})
	t.Equal(ret[3].Name, "A")
	t.Equal(ret[3].Index, []int{0, 0})
	t.Equal(ret[4].Name, "GetAllFieldT2")
	t.Equal(ret[5].Name, "C")
	t.Equal(ret[5].Index, []int{0, 2, 2})
	t.Equal(ret[6].Name, "D")
	t.Equal(ret[6].Index, []int{1, 1})

	ret = StructGetAllField(reflect.TypeOf(&GetAllFieldT5{}))
	t.Equal(len(ret), 2)

}

func TestStructGetAllFieldMap(ot *testing.T) {
	t := test.NewTestTools(ot)
	t1 := reflect.TypeOf(&GetAllFieldT1{})
	ret := StructGetAllFieldMap(t1)
	t.Equal(ret["A"].Index, []int{0, 0})
	t.Equal(ret["B"].Index, []int{2})
	t.Equal(ret["C"].Index, []int{0, 2, 2})
	t.Equal(ret["D"].Index, []int{1, 1})
	t.Equal(len(ret), 7)

	ret = StructGetAllFieldMap(reflect.TypeOf(&GetAllFieldT5{}))
	t.Equal(ret["A"].Index, []int{1})
	t.Equal(len(ret), 2)
}
