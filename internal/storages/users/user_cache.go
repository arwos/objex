package users

import (
	"sync"
)

type Cache struct {
	cache map[string]User
	mux   sync.RWMutex
}

func NewCache(count int) *Cache {
	return &Cache{
		cache: make(map[string]User, count),
	}
}

func (v *Cache) Get(login string) (User, bool) {
	v.mux.RLock()
	u, ok := v.cache[login]
	v.mux.RUnlock()
	return u, ok
}

func (v *Cache) Set(u User) {
	v.mux.Lock()
	v.cache[u.Login] = u
	v.mux.Unlock()
}

func (v *Cache) Delete(login string) {
	v.mux.Lock()
	delete(v.cache, login)
	v.mux.Unlock()
}

func (v *Cache) FlushAll() {
	v.mux.Lock()
	for login := range v.cache {
		delete(v.cache, login)
	}
	v.mux.Unlock()
}
