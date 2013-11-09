package sessionStore

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

type Manager struct {
	Provider
}

// load a store with that guid
// if that guid not exist will generate a new guid and a new store
func (manager *Manager) LoadStoreOrNewIfNotExist(guid string) (store *Store, err error) {
	// no need for lock manager,
	// if a store not exist ,it will get a new store with new random guid
	if manager.Exist(guid) {
		return manager.Get(guid)
	}
	return manager.New()
}
func (manager *Manager) New() (*Store, error) {
	guid, err := manager.generateGuid()
	if err != nil {
		return nil, err
	}
	return manager.Provider.NewByGuid(guid)
}

func (manager *Manager) generateGuid() (string, error) {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	guid := base64.URLEncoding.EncodeToString(b)
	return guid, nil
}
