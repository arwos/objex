package clients

import "io"

type Client interface {
	Connect(address string, opts ...any) (err error)
	Free() (n int64, err error)
	Close() (err error)

	Exist(filename string) (ok bool)
	Load(filename string) (r io.ReadCloser, err error)
	Save(filename string, r io.Reader) (err error)
	Delete(filename string) (err error)
	Size(filename string) (n int64, err error)
}
