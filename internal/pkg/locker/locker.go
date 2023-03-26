package locker

import (
	"sync"
)

type (
	locker struct {
		l    sync.Mutex
		data map[string]*mutex
	}
	mutex struct {
		l sync.Mutex
	}
	Locker interface {
		Mutex(name string) sync.Locker
	}
)

func New() Locker {
	return &locker{
		l:    sync.Mutex{},
		data: make(map[string]*mutex, 1000),
	}
}

func (v *locker) Mutex(name string) sync.Locker {
	v.l.Lock()
	mux, ok := v.data[name]
	if !ok {
		mux = new(mutex)
		v.data[name] = mux
	}
	v.l.Unlock()

	return mux
}

func (v *mutex) Lock() {
	v.l.Lock()
}

func (v *mutex) Unlock() {
	v.l.Unlock()
}
