package lib

import "github.com/bradfitz/gomemcache/memcache"

type Cacher struct {
	server string
	prefix string
	expires int32
	client *memcache.Client
}

func NewCacher(server string,prefix string,expires int32) *Cacher {
	return &Cacher{server:server,prefix:prefix,expires:expires, client:memcache.New(server)}
}

func (cache *Cacher)Exists(key string) bool {
	if _, err := cache.client.Get(cache.prefix+key); err != nil {
		return false
	}else {
		return true
	}
}

func (cache *Cacher)Get(key string) ([]byte, error) {
	if item, err := cache.client.Get(cache.prefix +key); err != nil {
		return nil, err
	}else {
		return item.Value, nil
	}
}

func (cache *Cacher)Put(key string, value []byte) error {
	item := &memcache.Item{Key:cache.prefix +key, Value:value}
	item.Expiration = cache.expires
	return cache.client.Set(item)
}

func (cache *Cacher)Delete(key string) error {
	return cache.client.Delete(key)
}

func (cache *Cacher)KeyCache(key string)string{
	return cache.prefix+key+""
}

func (cache *Cacher)KeyOrigin(key string){
	return cache.prefix+key
}
