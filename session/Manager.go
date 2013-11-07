package session

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

type Manager struct {
	Provider
}

func (manager *Manager) LoadStoreOrNewIfNotExist(guid string) (*Store, error) {
	// no need for lock manager,
	// if a store not exist ,it will get a new store with new random guid
	if manager.Exist(guid) {
		return manager.Get(guid), nil
	}
	return manager.NewStore()
}
func (manager *Manager) New() (*Store, error) {
	guid, err := manager.generateGuid()
	if err != nil {
		return nil, err
	}
	store := manager.Provider.NewByGuid(guid)
	return store, nil
}

func (manager *Manager) generateGuid() (string, error) {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}
	guid := base64.URLEncoding.EncodeToString(b)
	return guid, nil
}
