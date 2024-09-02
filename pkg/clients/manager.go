package clients

import "sync"

var (
	clients = map[string]func() Client{}
	mux     = sync.RWMutex{}
)

func Register(name string, initFunc func() Client) {
	mux.Lock()
	defer mux.Unlock()

	clients[name] = initFunc
}

func Get(name string) Client {
	mux.RLock()
	defer mux.RUnlock()

	if initFunc, ok := clients[name]; ok {
		return initFunc()
	}
	return nil
}
