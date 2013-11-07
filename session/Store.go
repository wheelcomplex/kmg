package session

import "sync"
//one session ,it is concurrent safe except GetAll
type Store struct{
	values map[string]interface {}
	guid string
	lock sync.RWMutex
}
func NewStore(guid string,values map[string]interface {})*Store{
	return &Store{guid:guid,values:values}
}
func (store *Store)Get(key string)(value interface {},exist bool){
	store.lock.RLock()
	defer store.lock.RUnlock()
	return store.values[key];
}
func (store *Store)Set(key string,value interface {}){
	store.lock.Lock()
	defer store.lock.Unlock()
	store.values[key] = value
}
func (store *Store)Delete(key string){
	store.lock.Lock()
	defer store.lock.Unlock()
	delete(store.values, key)
}
func (store *Store)Guid() string{
	return store.guid
}
func (store *Store)DeleteAll(){
	store.lock.Lock()
	defer store.lock.Unlock()
	store.values = make(map[string]interface {})
}
//return underlying value for save, it not concurrent safe
func (store *Store)GetAll() map[string]interface{}{
	return store.values
}
