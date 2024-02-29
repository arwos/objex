package ftp

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"time"

	ftplib "github.com/jlaffaye/ftp"
	"go.arwos.org/objex/internal/adapters/filestore/common"
	"go.osspkg.com/goppy/errors"
	"go.osspkg.com/goppy/xc"
)

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
	return v.connect(context.TODO(), func(c *ftplib.ServerConn) error {
		_, err := c.CurrentDir()
		return err
	})
}

func (v object) Del(ctx context.Context, filename string) error {
	return v.connect(ctx, func(c *ftplib.ServerConn) error {
		return c.Delete(filename)
	})
}

func (v object) Get(ctx context.Context, filename string, w io.Writer) error {
	return v.connect(ctx, func(c *ftplib.ServerConn) error {
		if _, err := c.FileSize(filename); err != nil {
			return fmt.Errorf("file not found [%s]", filename)
		}
		r, err := c.Retr(filename)
		if err != nil {
			return fmt.Errorf("open file [%s]", filename)
		}
		_, err = io.Copy(w, r)
		if err = errors.Wrap(err, r.Close()); err != nil {
			return fmt.Errorf("read file [%s]: %w", filename, err)
		}
		return nil
	})
}

func (v object) Set(ctx context.Context, filename string, r io.Reader) (common.THash, error) {
	dir := filepath.Dir(filename)

	dist, err := os.CreateTemp(os.TempDir(), "objex-*.bin")
	if err != nil {
		return "", err
	}
	tmpFileName := dist.Name()

	defer func() {
		dist.Close()           //nolint:errcheck
		os.Remove(tmpFileName) //nolint:errcheck
	}()

	hash := sha1.New()
	mw := io.MultiWriter(dist, hash)
	if _, err = io.Copy(mw, r); err != nil {
		return "", fmt.Errorf("write file [%s]: %w", tmpFileName, err)
	}

	_, err = dist.Seek(0, 0)
	if err != nil {
		return "", fmt.Errorf("seek file [%s]: %w", tmpFileName, err)
	}

	err = v.connect(ctx, func(c *ftplib.ServerConn) error {
		switch dir {
		case ".", "/":
		default:
			if err := c.MakeDir(dir); err != nil {
				return fmt.Errorf("create dir [%s]: %w", dir, err)
			}
		}
		err := c.Stor(filename, dist)
		if err = errors.Wrap(err, dist.Close()); err != nil {
			return fmt.Errorf("write file [%s]: %w", filename, err)
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	return common.THash(hex.EncodeToString(hash.Sum(nil))), nil
}

func (v object) connect(ctx context.Context, call func(c *ftplib.ServerConn) error) error {
	mctx, cncl := xc.Combine(v.ctx, ctx)
	defer cncl()

	host, login, passwd, dir, err := parseSetting(v.conf.Setting)
	if err != nil {
		return errors.Wrapf(err, "decode ftp setting")
	}
	conn, err := ftplib.Dial(host, ftplib.DialWithContext(mctx), ftplib.DialWithTimeout(15*time.Second))
	if err != nil {
		return errors.Wrapf(err, "open ftp connect")
	}
	defer conn.Quit() //nolint:errcheck
	if err = conn.Login(login, passwd); err != nil {
		return errors.Wrapf(err, "ftp authorization")
	}
	if len(dir) > 0 {
		if err = conn.ChangeDir(dir); err != nil {
			return errors.Wrapf(err, "change ftp dir")
		}
	}
	err = call(conn)
	return err
}

func parseSetting(uri string) (string, string, string, string, error) {
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
