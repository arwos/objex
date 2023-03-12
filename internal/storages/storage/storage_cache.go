package storage

import "sync"

type Cache struct {
	cache map[string]Store
	mux   sync.RWMutex
}

func NewCache(count int) *Cache {
	return &Cache{
		cache: make(map[string]Store, count),
	}
}

func (v *Cache) Get(name string) (Store, bool) {
	v.mux.RLock()
	s, ok := v.cache[name]
	v.mux.RUnlock()
	return s, ok
}

func (v *Cache) Set(s Store) {
	v.mux.Lock()
	v.cache[s.Name] = s
	v.mux.Unlock()
}

func (v *Cache) Delete(name string) {
	v.mux.Lock()
	delete(v.cache, name)
	v.mux.Unlock()
}

func (v *Cache) FlushAll() {
	v.mux.Lock()
	for name := range v.cache {
		delete(v.cache, name)
	}
	v.mux.Unlock()
}
