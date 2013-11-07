package ajkApi

import (
	"kmg/session"
)

// a lazy load session
type Session struct {
	//Self  *Peer
	//Other *Peer
	guid string
	store *session.Store
	manager *session.Manager
}
func NewSession(guid string,manager session.Manager)*Session{
	return &Session{guid:guid,manager:manager}
}
//lazy load session
func (sess *Session)GetStore()(store *Store,err error){
	err=sess.ConfirmSessionStart()
	if err!=nil{
		return
	}
	return sess.store,nil
}
func (sess *Session)GetGuid()string{
	return sess.guid
}
//确认session开始了,有store对象
func (sess *Session)ConfirmSessionStart()error{
	if sess.store!=nil{
		return nil
	}
	if sess.guid==""{
		sess.store,err:=sess.manager.New()
		sess.guid = sess.store.guid
		return err
	}
	sess.store,err:=sess.manager.LoadStoreOrNewIfNotExist(sess.guid)
	sess.guid = sess.store.guid
	return err
}
