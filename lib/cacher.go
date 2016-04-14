package lib

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"hash"
)

type Cacher struct {
	server  string
	prefix  string
	expires int32
	hasher  hash.Hash
	client  *memcache.Client
}

func NewCacher(server string, prefix string, expires int32) *Cacher {
	return &Cacher{
		server:  server,
		prefix:  prefix,
		expires: expires,
		hasher:  md5.New(),
		client:  memcache.New(server),
	}
}

func (cache *Cacher) Exists(key string) bool {
	if _, err := cache.client.Get(cache.getkey(key)); err != nil {
		return false
	} else {
		return true
	}
}

func (cache *Cacher) Get(key string) ([]byte, error) {
	if item, err := cache.client.Get(cache.getkey(key)); err != nil {
		return nil, err
	} else {
		return item.Value, nil
	}
}

func (cache *Cacher) Put(key string, value []byte) error {
	item := &memcache.Item{Key: cache.getkey(key), Value: value}
	item.Expiration = cache.expires
	return cache.client.Set(item)
}

func (cache *Cacher) Delete(key string) error {
	return cache.client.Delete(cache.getkey(key))
}

func (cache *Cacher) getkey(key string) string {
	cache.hasher.Reset()
	cache.hasher.Write([]byte(key))
	key = hex.EncodeToString(cache.hasher.Sum(nil))
	return cache.prefix + key
}
