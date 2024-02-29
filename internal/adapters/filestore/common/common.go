package common

import (
	"context"
	"io"
)

const (
	TYPE_LOCAL = "local"
	TYPE_FTP   = "ftp"
)

var typesList = map[string]struct{}{
	TYPE_LOCAL: {},
	TYPE_FTP:   {},
}

func IsValidType(s string) bool {
	_, ok := typesList[s]
	return ok
}

type THash string

type TProvider interface {
	Name() string
	Code() string
	Check() error

	Del(ctx context.Context, filename string) error
	Get(ctx context.Context, filename string, w io.Writer) error
	Set(ctx context.Context, filename string, r io.Reader) (THash, error)
}
