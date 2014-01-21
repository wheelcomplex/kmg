package memcacheProvider

import (
	"github.com/bradfitz/gomemcache/memcache"
)

type Provider struct {
	Client     *memcache.Client
	Prefix     string
	Expiration int32
}

func New(server ...string) *Provider {
	return &Provider{
		Client: memcache.New(server...),
	}
}
func (provider *Provider) Get(Id string) (Value []byte, Exist bool, err error) {
	item, err := provider.Client.Get(provider.Prefix + Id)
	if err == memcache.ErrCacheMiss {
		return nil, false, nil
	}
	Value = item.Value
	return
}
func (provider *Provider) Set(Id string, Value []byte) (err error) {
	item := &memcache.Item{
		Key:        provider.Prefix + Id,
		Value:      Value,
		Expiration: provider.Expiration,
	}
	err = provider.Client.Set(item)
	return
}
func (provider *Provider) Delete(Id string) (err error) {
	return provider.Client.Delete(provider.Prefix + Id)
}
