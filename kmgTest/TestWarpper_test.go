package kmgTest

//测试组件的测试

import "testing"

type TestTester struct {
	TestTools
	T1 int
}

func TestTestWarpper(ot *testing.T) {
	tester := &TestTester{
		T1: 1,
	}
	TestWarpper(ot, tester)

	if tester.T1 != 2 {
		ot.Errorf("[TestTestWarpper] tester.T1:%d!=2", tester.T1)
	}
}
func (t *TestTester) Test1() {
	t.Equal(1, 1)
	t.T1 = 2
}
