package kmgTest

//kmpTest Component
import "runtime"
import "fmt"
import "reflect"
import (
	"testing"
)

type Fatalfer interface {
	Fatalf(s string, v ...interface{})
}

type FatalferAware interface {
	SetFatalfer(T Fatalfer)
}

type TestTools struct {
	Fatalfer
}

func NewTestTools(T Fatalfer) *TestTools {
	return &TestTools{Fatalfer: T}
}

func (tools *TestTools) SetFatalfer(T Fatalfer) {
	tools.Fatalfer = T
}
func (tools *TestTools) Ok(expectTrue bool) {
	if !expectTrue {
		tools.assertFail("ok fail", 2)
	}
	return
}
func (tools *TestTools) Equal(get interface{}, expect interface{}) {
	if isEqual(expect, get) {
		return
	}
	msg := fmt.Sprintf("expect: %#v (%s) (%T)\nget: %#v (%s) (%T)", expect, expect, expect, get, get, get)
	if eGet, ok := get.(error); ok {
		msg += "\ngetError: " + eGet.Error()
	}
	tools.assertFail(msg, 2)
}
func (tools *TestTools) EqualMsg(get interface{}, expect interface{}, format string, args ...interface{}) {
	if isEqual(expect, get) {
		return
	}
	tools.assertFail(fmt.Sprintf("%s\nexpect:%#v (%T)\nget:%#v (%T)", fmt.Sprintf(format, args...), expect, expect, get, get), 2)
}
func (tools *TestTools) GetTestingT() *testing.T {
	return tools.Fatalfer.(*testing.T)
}

//an easy way to Printf(avoid of import fmt,and remove it ...)should for Debug only
func (tools *TestTools) Printf(format string, objs ...interface{}) (n int, err error) {
	return fmt.Printf(format, objs...)
}

//an easy way to Println (avoid of import fmt,and remove it ...)should for Debug only
func (tools *TestTools) Println(objs ...interface{}) (n int, err error) {
	return fmt.Println(objs...)
}

func (tools *TestTools) assertFail(msg string, skip int) {
	pc, file, line, ok := runtime.Caller(skip)
	line_info := ""
	if ok != false {
		line_info = fmt.Sprintf("%v:%v:%x", file, line, pc)
	}
	tools.Fatalf("%s\n%s", msg, line_info)
}

func isEqual(a interface{}, b interface{}) bool {
	if reflect.DeepEqual(a, b) {
		return true
	}
	rva := reflect.ValueOf(a)
	rvb := reflect.ValueOf(b)
	//every nil is same stuff...
	if isNil(rva) && isNil(rvb) {
		return true
	}
	return false
}

func isNil(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Interface, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
