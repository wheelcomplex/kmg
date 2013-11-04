package test

//测试组件的测试

import "testing"

type logFatalfer struct{
	call_num int
}

func (l *logFatalfer) Fatalf(s string, v...interface{}) {
	l.call_num+=1
}
func (l *logFatalfer) getCallNum() int {
	return l.call_num
}

func TestAssert(T *testing.T) {
	l := &logFatalfer{}
	Assert(l, true, false)
	if l.getCallNum() != 1 {
		T.Fatal("fail!")
	}
	Assert(l, true, true)
	if l.getCallNum() != 1 {
		T.Fatal("fail!")
	}
	Assert(l, false, false)
	if l.getCallNum() != 1 {
		T.Fatal("fail!")
	}
}
func TestIsEqual(T *testing.T) {
	Assert(T, true, isEqual(map[string]interface{}{"a":1}, map[string]interface{}{"a":1}))
}
