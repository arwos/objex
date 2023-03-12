package providers

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/deweppro/go-sdk/errors"
	"github.com/deweppro/go-sdk/file"
	"github.com/deweppro/go-static"
	"github.com/deweppro/goppy/plugins/web"
)

type LocalProvider struct {
	conf Item
	ctx  context.Context
}

func NewLocalProvider(ctx context.Context, c Item) *LocalProvider {
	return &LocalProvider{
		ctx:  ctx,
		conf: c,
	}
}

func (v *LocalProvider) Name() string {
	return v.conf.Name
}

func (v *LocalProvider) Code() string {
	return v.conf.Code
}

func (v *LocalProvider) Check() error {
	if !file.Exist(v.conf.Setting) {
		return os.MkdirAll(v.conf.Setting, 0755)
	}
	return nil
}

func (v *LocalProvider) GetFile(filename string, ctx web.Context) {
	origFile := filepath.Join(v.conf.Setting, filename)
	if !file.Exist(origFile) {
		ctx.Error(http.StatusNotFound, fmt.Errorf("file not found"))
		return
	}

	if len(v.conf.Prefix) != 0 {
		ctx.Redirect(v.conf.Prefix + "/" + filename)
		return
	}

	resp, err := os.Open(origFile)
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err)
		return
	}
	defer resp.Close() //nolint:errcheck

	contentType := static.DetectContentType(filename, nil)
	ctx.Header().Set("Content-Type", contentType)
	ctx.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filename)))
	ctx.Response().WriteHeader(http.StatusOK)
	io.Copy(ctx.Response(), resp) //nolint:errcheck
}

func (v *LocalProvider) SaveFile(filename string, r io.ReadCloser) (string, error) {
	origFile := filepath.Join(v.conf.Setting, filename)
	if file.Exist(origFile) {
		return "", fmt.Errorf("file alredy exist")
	}

	err := os.MkdirAll(filepath.Dir(origFile), 0744)
	if err != nil {
		return "", err
	}

	dist, err := os.OpenFile(origFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}

	hash := sha1.New()
	mw := io.MultiWriter(dist, hash)
	_, err = io.Copy(mw, r)

	err = errors.Wrap(err, r.Close(), dist.Close())
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (v *LocalProvider) DeleteFile(filename string) error {
	origFile := filepath.Join(v.conf.Setting, filename)
	if file.Exist(origFile) {
		return nil
	}

	return os.Remove(origFile)
}
