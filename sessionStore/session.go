package sessionStore

import (
	"labix.org/v2/mgo/bson"
)

//any call in this object not thread safe
type Session struct {
	Id      string
	data    map[string]interface{}
	manager *Manager
}

func (sess *Session) GetIntOrZero(key string) (i int) {
	value, ok := sess.data[key]
	if !ok {
		return 0
	}
	i, ok = value.(int)
	if !ok {
		return 0
	}
	return
}
func (sess *Session) Set(key string, value interface{}) {
	sess.data[key] = value
}

func (sess *Session) Get(key string) (value interface{}, ok bool) {
	value, ok = sess.data[key]
	return
}

//delete current session,delete all data in this session ,create a new one,
//it will panic if an error happened
func (sess *Session) DeleteAndNewSession() {
	err := sess.manager.Provider.Delete(sess.Id)
	if err != nil {
		panic(err)
	}
	newSess, err := sess.manager.newSession()
	if err != nil {
		panic(err)
	}
	*sess = *newSess
}

//delete current session,and delete all data in this session
func (sess *Session) DeleteSession() {
	panic("[Session.DeleteSession]not implement!")
}

func newSession(Manager *Manager, Id string) *Session {
	return &Session{
		Id:      Id,
		data:    make(map[string]interface{}),
		manager: Manager,
	}
}

func unmarshalSession(value []byte, Manager *Manager, Id string) (session *Session, err error) {
	session = &Session{
		Id:      Id,
		data:    make(map[string]interface{}),
		manager: Manager,
	}
	err = bson.Unmarshal(value, &session.data)
	if err != nil {
		return
	}
	return
}

func (session *Session) marshal() (out []byte, err error) {
	return bson.Marshal(session.data)
}
