package providers

import (
	"fmt"
	"io"
	"net/url"

	"github.com/deweppro/goppy/plugins/web"
)

const (
	TypeLocal = "local"
	TypeFTP   = "ftp"
)

var typesList = map[string]struct{}{
	TypeLocal: {},
	TypeFTP:   {},
}

type Provider interface {
	Name() string
	Code() string
	Check() error
	GetFile(filename string, ctx web.Context)
	SaveFile(filename string, r io.ReadCloser) (string, error)
	DeleteFile(filename string) error
}

func IsValidType(s string) bool {
	_, ok := typesList[s]
	return ok
}

func ParseFTP(uri string) (string, string, string, string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", "", "", "", err
	}
	host := u.Host
	login := u.User.Username()
	passwd, ok := u.User.Password()
	if !ok {
		return "", "", "", "", fmt.Errorf("passwd is empty")
	}
	dir := u.Path

	return host, login, passwd, dir, nil
}
