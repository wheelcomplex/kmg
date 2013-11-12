package kmgReflect

import (
	"github.com/bronze1man/kmg/test"
	"reflect"
	"testing"
)

type ta struct {
}

func TestGetFullName(ot *testing.T) {
	t := test.NewTestTools(ot)
	name, ok := GetTypeFullName(reflect.TypeOf(""))
	t.Equal(name, "string")
	t.Equal(ok, true)

	name, ok = GetTypeFullName(reflect.TypeOf(1))
	t.Equal(name, "int")
	t.Equal(ok, true)

	name, ok = GetTypeFullName(reflect.TypeOf(&ta{}))
	t.Equal(name, "github.com/bronze1man/kmg/kmgReflect.ta")
	t.Equal(ok, true)

	name, ok = GetTypeFullName(reflect.TypeOf([]string{}))
	t.Equal(name, "")
	t.Equal(ok, false)

}
