package local

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.arwos.org/objex/internal/adapters/filestore/common"
	"go.osspkg.com/goppy/errors"
	"go.osspkg.com/goppy/iofile"
)

const perm os.FileMode = 0744

type object struct {
	conf common.ConfigItem
	ctx  context.Context
}

func New(ctx context.Context, cfg common.ConfigItem) common.TProvider {
	return &object{
		ctx:  ctx,
		conf: cfg,
	}
}

func (v object) Name() string {
	return v.conf.Name
}

func (v object) Code() string {
	return v.conf.Code
}

func (v object) Check() error {
	if !iofile.Exist(v.conf.Setting) {
		return os.MkdirAll(v.conf.Setting, perm)
	}
	return nil
}

func (v object) Del(_ context.Context, filename string) error {
	path := v.buildPath(filename)
	if !iofile.Exist(path) {
		return nil
	}
	return os.Remove(path)
}

func (v object) Get(_ context.Context, filename string, w io.Writer) error {
	path := v.buildPath(filename)

	if !iofile.Exist(path) {
		return fmt.Errorf("file not found [%s]", path)
	}

	r, err := os.OpenFile(path, os.O_RDONLY, perm)
	if err != nil {
		return fmt.Errorf("open file [%s]: %w", path, err)
	}

	_, err = io.Copy(w, r)
	if err = errors.Wrap(err, r.Close()); err != nil {
		return fmt.Errorf("read file [%s]: %w", path, err)
	}
	return nil
}

func (v object) Set(_ context.Context, filename string, r io.Reader) (common.THash, error) {
	path := v.buildPath(filename)

	err := os.MkdirAll(filepath.Dir(path), perm)
	if err != nil {
		return "", fmt.Errorf("create folder: %w", err)
	}

	dist, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return "", fmt.Errorf("open file [%s]: %w", path, err)
	}

	hash := sha1.New()
	mw := io.MultiWriter(dist, hash)
	_, err = io.Copy(mw, r)

	if err = errors.Wrap(err, dist.Close()); err != nil {
		return "", fmt.Errorf("write file [%s]: %w", path, err)
	}

	return common.THash(hex.EncodeToString(hash.Sum(nil))), nil
}

func (v object) buildPath(filename string) string {
	return filepath.Join(v.conf.Setting, filename)
}
