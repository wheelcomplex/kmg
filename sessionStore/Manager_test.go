package sessionStore

import "testing"
import "github.com/bronze1man/kmg/kmgTest"

func TestManager(ot *testing.T) {
	t := kmgTest.NewTestTools(ot)
	manager := &Manager{NewMemoryProvider()}
	store, err := manager.LoadStoreOrNewIfNotExist("1")
	t.Equal(err, nil)
	store.Set("A", 5)

	store, err = manager.LoadStoreOrNewIfNotExist(store.Guid())
	t.Equal(err, nil)
	value, ok := store.Get("A")
	t.Equal(ok, true)
	t.Equal(value, 5)
}
