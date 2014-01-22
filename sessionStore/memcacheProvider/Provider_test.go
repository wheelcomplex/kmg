package memcacheProvider

import (
	"github.com/bronze1man/kmg/kmgTest"
	"os"
	"testing"
)

func Test(ot *testing.T) {
	kmgTest.TestWarpper(ot, &Tester{})
}

type Tester struct {
	kmgTest.TestTools
}

func (t *Tester) TestProvider() {
	host := os.Getenv("TEST_MEMCACHE_HOST")
	if host == "" {
		t.GetTestingT().Skip(`need memcache host
export TEST_MEMCACHE_HOST=127.0.0.1:11211

skip this test ...`)
		return
	}
	p := New(host)
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
