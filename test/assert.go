package test

//kmpTest Component
import "runtime"
import "fmt"
import "reflect"

type Fatalfer interface {
	Fatalf(s string, v ...interface{})
}

type FatalferAware interface {
	SetFatalfer(T Fatalfer)
}

type TestTools struct {
	T Fatalfer
}

func NewTestTools(T Fatalfer) *TestTools {
	return &TestTools{T: T}
}
func (tools *TestTools) SetFatalfer(T Fatalfer) {
	tools.T = T
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
	msg := fmt.Sprintf("expect: %#v (%T)\nget: %#v (%T)", expect, expect, get, get)
	if eGet, ok := get.(error); ok {
		msg += "\ngetError: " + eGet.Error()
	}
	tools.assertFail(msg, 2)
}
func (tools *TestTools) EqualMsg(get interface{}, expect interface{}, msg string) {
	if isEqual(expect, get) {
		return
	}
	tools.assertFail(fmt.Sprintf("%s\nexpect:%#v (%T)\nget:%#v (%T)", msg, expect, expect, get, get), 2)
}

func (tools *TestTools) assertFail(msg string, skip int) {
	pc, file, line, ok := runtime.Caller(skip)
	line_info := ""
	if ok != false {
		line_info = fmt.Sprintf("%v:%v:%x", file, line, pc)
	}
	tools.T.Fatalf("%s\n%s", msg, line_info)
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
