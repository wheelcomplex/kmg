package test

//测试组件的测试

import "testing"

type logFatalfer struct {
	call_num int
}

func (l *logFatalfer) Fatalf(s string, v ...interface{}) {
	l.call_num += 1
}
func (l *logFatalfer) getCallNum() int {
	return l.call_num
}

func TestEqual(ot *testing.T) {
	l := &logFatalfer{}
	t := NewTestTools(l)
	t.Equal(true, false)
	if l.getCallNum() != 1 {
		ot.Fatal("fail!")
	}
	t.Equal(true, true)
	if l.getCallNum() != 1 {
		ot.Fatal("fail!")
	}
	t.Equal(false, false)
	if l.getCallNum() != 1 {
		ot.Fatal("fail!")
	}
}
func TestIsEqual(ot *testing.T) {
	t := NewTestTools(ot)
	t.Equal(true, isEqual(map[string]interface{}{"a": 1}, map[string]interface{}{"a": 1}))
}
