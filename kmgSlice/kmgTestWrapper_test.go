package kmgSlice

import "testing"
import "github.com/bronze1man/kmg/kmgTest"

type Tester struct {
	kmgTest.TestTools
}

func Test1(t *testing.T) {
	kmgTest.TestWarpper(t, &Tester{})
}
