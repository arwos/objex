package users

import (
	"sync"
)

type Cache struct {
	cache  map[string]User
	ids    map[uint64]string
	tokens map[string]uint64
	mux    sync.RWMutex
}

func NewCache(count int) *Cache {
	return &Cache{
		cache:  make(map[string]User, count),
		ids:    make(map[uint64]string, count),
		tokens: make(map[string]uint64, count),
	}
}

/**********************************************************************************************************************/

func (v *Cache) getByLogin(login string) (User, bool) {
	u, ok := v.cache[login]
	return u, ok
}

func (v *Cache) getByID(uid uint64) (User, bool) {
	login, ok := v.ids[uid]
	if !ok {
		return User{}, false
	}
	return v.getByLogin(login)
}

func (v *Cache) getByToken(token string) (User, bool) {
	uid, ok := v.tokens[token]
	if !ok {
		return User{}, false
	}
	return v.getByID(uid)
}

/**********************************************************************************************************************/

func (v *Cache) GetByLogin(login string) (User, bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()
	return v.getByLogin(login)
}

func (v *Cache) GetByID(uid uint64) (User, bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()
	return v.getByID(uid)
}

func (v *Cache) GetByToken(token string) (User, bool) {
	v.mux.RLock()
	defer v.mux.RUnlock()
	return v.getByToken(token)
}

func (v *Cache) Setup(u User, tokens ...string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	v.cache[u.Login] = u
	v.ids[u.ID] = u.Login
	for _, token := range tokens {
		v.tokens[token] = u.ID
	}
}

func (v *Cache) SetToken(uid uint64, tokens ...string) bool {
	v.mux.Lock()
	defer v.mux.Unlock()

	if _, ok := v.ids[uid]; !ok {
		return false
	}
	for _, token := range tokens {
		v.tokens[token] = uid
	}
	return true
}

func (v *Cache) DeleteToken(uid uint64, token string) bool {
	v.mux.Lock()
	defer v.mux.Unlock()

	if _, ok := v.ids[uid]; !ok {
		return false
	}
	if id, ok := v.tokens[token]; ok && uid == id {
		delete(v.tokens, token)
		return true
	}
	return false
}

func (v *Cache) Delete(login string) {
	v.mux.Lock()
	defer v.mux.Unlock()

	u, ok := v.cache[login]
	if !ok {
		return
	}
	delete(v.ids, u.ID)
	delete(v.cache, u.Login)
	for token, id := range v.tokens {
		if id == u.ID {
			delete(v.tokens, token)
		}
	}
}

func (v *Cache) FlushAll() {
	v.mux.Lock()
	defer v.mux.Unlock()

	for key := range v.cache {
		delete(v.cache, key)
	}
	for key := range v.ids {
		delete(v.ids, key)
	}
	for key := range v.tokens {
		delete(v.tokens, key)
	}
}
