package kmgGob

import (
	"github.com/bronze1man/kmg/kmgTest"
	"testing"
)

func Test(t *testing.T) {
	kmgTest.TestWarpper(t, &S{})
}

type S struct {
	kmgTest.TestTools
}
