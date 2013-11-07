package session

type Provider interface{
	NewByGuid(guid string)(*Store,error)
	//read a store,if this store not exist,return a new empty store
	Read(guid string) (*Store,error)
	//save a store to someplace,should not change this store after save.
	Save(session *Store) error
	Exist(guid string) bool
	//delete a store
	Delete(guid string) error
	GC()
}

