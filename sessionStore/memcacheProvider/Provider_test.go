package memcacheProvider

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

func (t *Tester) TestProvider() {
	p := New("10.1.1.8:11211")
	p.Prefix = "memcache_test_"
	err := p.Set("g881r0H-B4fIGF8ktUWeUg==", []byte("2"))
	t.Equal(err, nil)

	v, ok, err := p.Get("g881r0H-B4fIGF8ktUWeUg==")
	t.Equal(err, nil)
	t.Equal(ok, true)
	t.Equal(v, []byte("2"))

	err = p.Delete("g881r0H-B4fIGF8ktUWeUg==")
	t.Equal(err, nil)

	v, ok, err = p.Get("g881r0H-B4fIGF8ktUWeUg==")
	t.Equal(err, nil)
	t.Equal(ok, false)

	err = p.Set("g881r0H-B4fIGF8ktUWeUg==", []byte("3"))
	t.Equal(err, nil)
	v, ok, err = p.Get("g881r0H-B4fIGF8ktUWeUg==")
	t.Equal(err, nil)
	t.Equal(ok, true)
	t.Equal(v, []byte("3"))
}
