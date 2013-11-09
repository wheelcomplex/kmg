package sessionStore

import "testing"
import "kmg/test"

func TestStore(ot *testing.T) {
	t := test.NewTestTools(ot)
	store := NewStore("1", make(map[string]interface{}))
	value, ok := store.Get("A")
	t.Equal(ok, false)

	store.Set("A", 1)
	value, ok = store.Get("A")
	t.Equal(ok, true)
	t.Equal(value, 1)

	store.Delete("A")

	value, ok = store.Get("A")
	t.Equal(ok, false)

	t.Equal(store.Guid(), "1")

	store.Set("A", 1)
	store.DeleteAll()

	value, ok = store.Get("A")
	t.Equal(ok, false)
}
