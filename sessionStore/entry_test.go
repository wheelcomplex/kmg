package sessionStore

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func Test(ot *testing.T) {
	kmgTest.TestWarpper(ot, &Tester{})
}

type Tester struct {
	kmgTest.TestTools
}
