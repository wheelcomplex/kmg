package kmgTest

import "fmt"
import "reflect"
import (
	"github.com/bronze1man/kmg/kmgDebug" //TODO 移除这个依赖?
	"testing"
)

type TestingTB interface {
	FailNow()
}

type TestingTBAware interface {
	SetTestingTB(T TestingTB)
}

type TestTools struct {
	TestingTB
}

func NewTestTools(T TestingTB) *TestTools {
	return &TestTools{TestingTB: T}
}

func (tools *TestTools) SetTestingTB(T TestingTB) {
	tools.TestingTB = T
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
	msg := fmt.Sprintf("\texpect2: %#v (%s) (%T)\n\tget1: %#v (%s) (%T)", expect, expect, expect, get, get, get)
	if eGet, ok := get.(error); ok {
		msg += "\ngetError: " + eGet.Error()
	}
	tools.assertFail(msg, 2)
}
func (tools *TestTools) EqualMsg(get interface{}, expect interface{}, format string, args ...interface{}) {
	if isEqual(expect, get) {
		return
	}
	tools.assertFail(fmt.Sprintf(`%s
expect2:%#v (%T)
get1:%#v (%T)`, fmt.Sprintf(format, args...), expect, expect, get, get), 2)
}
func (tools *TestTools) GetTestingT() *testing.T {
	return tools.TestingTB.(*testing.T)
}

//an easy way to Printf(avoid of import fmt,and remove it ...)should for Debug only
func (tools *TestTools) Printf(format string, objs ...interface{}) (n int, err error) {
	return fmt.Printf(format, objs...)
}

//an easy way to Println (avoid of import fmt,and remove it ...)should for Debug only
func (tools *TestTools) Println(objs ...interface{}) (n int, err error) {
	return fmt.Println(objs...)
}
func (tools *TestTools) Fatalf(format string, objs ...interface{}) {
	tools.Printf("\n%s\n\n%s\n", fmt.Sprintf(format, objs...),
		kmgDebug.GetCurrentStack(1).ToString())
	tools.TestingTB.FailNow()
}
func (tools *TestTools) assertFail(msg string, skip int) {
	/*
		pc, file, line, ok := runtime.Caller(skip)
		line_info := ""
		if ok != false {
			line_info = fmt.Sprintf("%v:%v:%x", file, line, pc)
		}
	*/
	tools.Printf(`----------------------------------
%s

%s----------------------------------
`, msg, kmgDebug.GetCurrentStack(skip).ToString())
	tools.TestingTB.FailNow()
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
