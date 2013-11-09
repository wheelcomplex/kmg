package sessionStore

import "testing"
import "kmg/test"

func TestManagerProvider(ot *testing.T) {
	t := test.NewTestTools(ot)
	provider := NewMemoryProvider()

	ret := provider.Exist("1")
	t.Equal(ret, false)

	store, err := provider.NewByGuid("1")
	t.Equal(err, nil)
	store.Set("A", 5)
	value, ok := store.Get("A")
	t.Equal(ok, true)
	t.Equal(value, 5)

	ret = provider.Exist("1")
	t.Equal(ret, true)

	store, err = provider.Get("1")
	t.Equal(err, nil)
	value, ok = store.Get("A")
	t.Equal(ok, true)
	t.Equal(value, 5)

}
