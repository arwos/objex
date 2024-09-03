package s3

import "go.arwos.org/objex/pkg/clients"

func init() {
	clients.Register("s3", func() clients.Client { return &Client{} })
}
