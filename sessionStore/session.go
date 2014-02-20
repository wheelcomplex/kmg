package sessionStore

import (
	"fmt"
	"github.com/bronze1man/kmg/typeTransform"
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
func (sess *Session) GetBool(key string) bool {
	value, ok := sess.data[key]
	if !ok {
		return false
	}
	b, ok := value.(bool)
	if !ok {
		return false
	}
	return b
}
func (sess *Session) Set(key string, value interface{}) {
	sess.data[key] = value
}

func (sess *Session) Get(key string) (value interface{}, ok bool) {
	value, ok = sess.data[key]
	return
}

func (sess *Session) GetWithType(key string, out interface{}) (err error) {
	value, ok := sess.data[key]
	if !ok {
		return fmt.Errorf("[GetWithType]key:%s not found", key)
	}
	return typeTransform.Transform(value, out)
}

//delete current session,delete all data in this session ,create a new one,
//it will panic if an error happened
func (sess *Session) DeleteAndNewSession() (err error) {
	err = sess.manager.Provider.Delete(sess.Id)
	if err != nil {
		return
	}
	newSess, err := sess.manager.newSession()
	if err != nil {
		return
	}
	*sess = *newSess
	return
}

//delete current session,and delete all data in this session
func (sess *Session) DeleteSession() (err error) {
	panic("[Session.DeleteSession]not implement!")
}

//save this Session and reload it from SessionProvider,should only used for test
func (sess *Session) SaveAndReload() (err error) {
	err = sess.manager.Save(sess)
	if err != nil {
		return
	}
	newSess, err := sess.manager.Load(sess.Id)
	if err != nil {
		return
	}
	*sess = *newSess
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
