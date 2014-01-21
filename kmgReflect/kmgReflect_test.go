package kmgReflect

import (
	"github.com/bronze1man/kmg/kmgTest"
	"reflect"
	"testing"
)

type ta struct {
}

func TestGetFullName(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	name := GetTypeFullName(reflect.TypeOf(""))
	t.Equal(name, "string")

	name = GetTypeFullName(reflect.TypeOf(1))
	t.Equal(name, "int")

	name = GetTypeFullName(reflect.TypeOf(&ta{}))
	t.Equal(name, "github.com/bronze1man/kmg/kmgReflect.ta")

	name = GetTypeFullName(reflect.TypeOf([]string{}))
	t.Equal(name, "")

}
