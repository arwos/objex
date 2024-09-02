package filesystem

import "go.arwos.org/objex/pkg/clients"

func init() {
	clients.Register("fs", func() clients.Client { return &Client{} })
}
