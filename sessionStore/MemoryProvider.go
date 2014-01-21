package sessionStore

import (
	"sync"
)

type MemoryProvider struct {
	data map[string][]byte
	lock sync.RWMutex
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{data: make(map[string][]byte)}
}
func (provider *MemoryProvider) Get(Id string) (Value []byte, Exist bool, err error) {
	provider.lock.RLock()
	defer provider.lock.RUnlock()
	Value, Exist = provider.data[Id]
	return
}
func (provider *MemoryProvider) Set(Id string, Value []byte) (err error) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	provider.data[Id] = Value
	return
}
func (provider *MemoryProvider) Delete(Id string) (err error) {
	provider.lock.Lock()
	defer provider.lock.Unlock()
	delete(provider.data, Id)
	return
}

/*
//do not use &MemoryProvider{} init
type MemoryProvider struct {
	stores map[string]*Store
	lock   sync.RWMutex
}

func NewMemoryProvider() *MemoryProvider {
	return &MemoryProvider{stores: make(map[string]*Store)}
}

var GuidNotExistErr = errors.New("guid not exist")

func (memoryProvider *MemoryProvider) NewByGuid(guid string) (store *Store, err error) {
	memoryProvider.lock.Lock()
	defer memoryProvider.lock.Unlock()
	store = NewStore(guid, make(map[string]interface{}))
	memoryProvider.stores[guid] = store
	return store, nil
}

func (memoryProvider *MemoryProvider) Get(guid string) (store *Store, err error) {
	memoryProvider.lock.RLock()
	defer memoryProvider.lock.RUnlock()
	store, ok := memoryProvider.stores[guid]
	if ok {
		return store, nil
	}
	return nil, GuidNotExistErr
}

func (memoryProvider *MemoryProvider) Exist(guid string) bool {
	memoryProvider.lock.RLock()
	defer memoryProvider.lock.RUnlock()
	_, ok := memoryProvider.stores[guid]
	return ok
}

func (memoryProvider *MemoryProvider) Delete(guid string) error {
	memoryProvider.lock.Lock()
	defer memoryProvider.lock.Unlock()
	delete(memoryProvider.stores, guid)
	return nil
}

func (memoryProvider *MemoryProvider) Save(store *Store) error {
	return nil
}

//TODO
func (memoryProvider *MemoryProvider) GC() {

}
*/
