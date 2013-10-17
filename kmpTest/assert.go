package kmpTest

//kmpTest Component
import "runtime"
import "fmt"
import "reflect"

type Fatalfer interface{
	Fatalf(s string, v...interface{})
}

func Assert(T Fatalfer, get interface{}, expect interface{}) {
	if isEqual(expect, get) {
		return
	}
	assertFail(T, fmt.Sprintf("expect:%#v\nget:%#v", expect, get), 2)
}
func AssertMsg(T Fatalfer, get interface{}, expect interface{}, msg string) {
	if isEqual(expect, get) {
		return
	}
	assertFail(T, fmt.Sprintf("%s\nexpect:%#v\nget:%#v", msg, expect, get), 2)
}
func isEqual(a interface{}, b interface{}) bool {
	if reflect.DeepEqual(a, b) {
		return true
	}
	rva := reflect.ValueOf(a);
	rvb := reflect.ValueOf(b);
	//every nil is same stuff...
	if isNil(rva) && isNil(rvb) {
		return true
	}
	return false
}
func assertFail(T Fatalfer, msg string, skip int) {
	pc, file, line, ok := runtime.Caller(skip)
	line_info := ""
	if ok != false {
		line_info = fmt.Sprintf("%v:%v:%x", file, line, pc)
	}
	T.Fatalf("%s\n%s", msg, line_info)
}
func isNil(rv reflect.Value) bool {
	switch rv.Kind(){
	case reflect.Invalid:
		return true
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Interface, reflect.Ptr, reflect.Slice:
		return rv.IsNil()
	default:
		return false
	}
}
