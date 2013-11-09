package ajkApi

import (
	"github.com/bronze1man/kmg/sessionStore"
)

// a lazy load session
type Session struct {
	//Self  *Peer
	//Other *Peer
	guid         string
	store        *sessionStore.Store
	storeManager *sessionStore.Manager
}

func NewSession(guid string, storeManager *sessionStore.Manager) *Session {
	return &Session{guid: guid, storeManager: storeManager}
}

//lazy load session
func (sess *Session) GetStore() (store *sessionStore.Store, err error) {
	err = sess.ConfirmSessionStart()
	if err != nil {
		return nil, err
	}
	return sess.store, nil
}
func (sess *Session) GetGuid() string {
	return sess.guid
}

//确认session开始了,有store对象
func (sess *Session) ConfirmSessionStart() (err error) {
	if sess.store != nil {
		return nil
	}
	if sess.guid == "" {
		sess.store, err = sess.storeManager.New()
		sess.guid = sess.store.Guid()
		return err
	}
	sess.store, err = sess.storeManager.LoadStoreOrNewIfNotExist(sess.guid)
	sess.guid = sess.store.Guid()
	return err
}
