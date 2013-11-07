package session

import (
	"sync"
	"errors"
)

type MemoryProvider struct{
	stores map[string]*Store
	lock sync.RWMutex
}
var GuidNotExistErr = errors.New("guid not exist")

func (memoryProvider *MemoryProvider)NewByGuid(guid string)(store *Store,error){
	memoryProvider.lock.Lock()
	defer memoryProvider.lock.Unlock()
	store:= NewStore(guid,make(map[string]interface {}))
	return store,nil
}

func (memoryProvider *MemoryProvider)Read(guid string)(store *Store,error){
	memoryProvider.lock.RLock()
	defer memoryProvider.lock.RUnlock()
	store,ok := memoryProvider.stores[guid]
	if ok{
		return store,nil
	}
	return nil,GuidNotExistErr
}

func (memoryProvider *MemoryProvider)Exist(guid string)bool{
	memoryProvider.lock.RLock()
	defer memoryProvider.lock.RUnlock()
	_,ok := memoryProvider.stores[guid]
	return ok
}

func (memoryProvider *MemoryProvider)Delete(guid string)error{
	memoryProvider.lock.Lock()
	defer memoryProvider.lock.Unlock()
	delete(memoryProvider.stores , guid)
	return nil
}

func (memoryProvider *MemoryProvider)Save(store *Store)error{
	return nil
}

//TODO
func (memoryProvider *MemoryProvider)GC(){

}


