package memcacheProvider

import (
	"github.com/bradfitz/gomemcache/memcache"
	"net"
)

type Provider struct {
	Client     *memcache.Client
	Prefix     string
	Expiration int32
}

func New(server ...string) *Provider {
	if len(server) == 0 {
		panic("[memcacheProvider.New] len(server)==0")
	}
	return &Provider{
		Client: memcache.New(server...),
	}
}
func (provider *Provider) Get(Id string) (Value []byte, Exist bool, err error) {
	//workaround for memcache network error
	for i := 0; i < 2; i++ {
		Value, Exist, err = provider.get(Id)
		if err == nil {
			return
		}
		_, ok := err.(net.Error)
		if !ok {
			return
		}
	}
	return
}
func (provider *Provider) get(Id string) (Value []byte, Exist bool, err error) {
	item, err := provider.Client.Get(provider.Prefix + Id)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return nil, false, nil
		}
		return
	}
	Value = item.Value
	Exist = true
	return
}
func (provider *Provider) Set(Id string, Value []byte) (err error) {
	//workaround for memcache network error
	for i := 0; i < 2; i++ {
		err = provider.set(Id, Value)
		if err == nil {
			return
		}
		_, ok := err.(net.Error)
		if !ok {
			return
		}
	}
	return
}
func (provider *Provider) set(Id string, Value []byte) (err error) {
	item := &memcache.Item{
		Key:        provider.Prefix + Id,
		Value:      Value,
		Expiration: provider.Expiration,
	}
	err = provider.Client.Set(item)
	return
}
func (provider *Provider) Delete(Id string) (err error) {
	//workaround for memcache network error
	for i := 0; i < 2; i++ {
		err = provider.delete(Id)
		if err == nil {
			return
		}
		_, ok := err.(net.Error)
		if !ok {
			return
		}
	}
	return
}

func (provider *Provider) delete(Id string) (err error) {
	err = provider.Client.Delete(provider.Prefix + Id)
	if err == nil {
		return nil
	}
	if err == memcache.ErrCacheMiss {
		return nil
	}
	return err
}
