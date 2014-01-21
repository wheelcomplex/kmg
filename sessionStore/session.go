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
