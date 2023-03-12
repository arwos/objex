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
	"time"

	"github.com/deweppro/go-sdk/errors"
	"github.com/deweppro/go-static"
	"github.com/deweppro/goppy/plugins/web"
	"github.com/jlaffaye/ftp"
)

type FTPProvider struct {
	conf Item
	ctx  context.Context
}

func NewFTPProvider(ctx context.Context, c Item) *FTPProvider {
	return &FTPProvider{
		ctx:  ctx,
		conf: c,
	}
}

func (v *FTPProvider) connect(call func(*ftp.ServerConn) error) error {
	host, login, passwd, dir, err := ParseFTP(v.conf.Setting)
	if err != nil {
		return errors.Wrapf(err, "decode ftp setting")
	}
	c, err := ftp.Dial(host, ftp.DialWithContext(v.ctx), ftp.DialWithTimeout(15*time.Second))
	if err != nil {
		return errors.Wrapf(err, "open ftp connect")
	}
	defer c.Quit() //nolint:errcheck
	if err = c.Login(login, passwd); err != nil {
		return errors.Wrapf(err, "ftp authorization")
	}
	if len(dir) != 0 {
		if err = c.ChangeDir(dir); err != nil {
			return errors.Wrapf(err, "change ftp dir")
		}
	}
	err = call(c)
	return err
}

func (v *FTPProvider) Name() string {
	return v.conf.Name
}

func (v *FTPProvider) Code() string {
	return v.conf.Code
}

func (v *FTPProvider) Check() error {
	return v.connect(func(conn *ftp.ServerConn) error {
		_, err := conn.CurrentDir()
		return err
	})
}

func (v *FTPProvider) GetFile(filename string, ctx web.Context) {
	err := v.connect(func(conn *ftp.ServerConn) error {
		_, err := conn.FileSize(filename)
		if err != nil {
			return fmt.Errorf("file not found")
		}
		if len(v.conf.Prefix) != 0 {
			ctx.Redirect(v.conf.Prefix + "/" + filename)
			return nil
		}

		resp, err := conn.Retr(filename)
		if err != nil {
			return err
		}
		defer resp.Close() //nolint:errcheck

		contentType := static.DetectContentType(filename, nil)
		ctx.Header().Set("Content-Type", contentType)
		ctx.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filename)))
		ctx.Response().WriteHeader(http.StatusOK)
		_, err = io.Copy(ctx.Response(), resp)
		return err
	})
	if err != nil {
		ctx.Error(http.StatusInternalServerError, err)
	}
}

func (v *FTPProvider) SaveFile(filename string, r io.ReadCloser) (string, error) {
	dir := filepath.Dir(filename)
	hash := sha1.New()
	tf, err := os.CreateTemp(os.TempDir(), "artifactory-*.bin")
	if err != nil {
		return "", err
	}
	defer func() {
		tfName := tf.Name()
		tf.Close()        //nolint:errcheck
		os.Remove(tfName) //nolint:errcheck
	}()
	mw := io.MultiWriter(tf, hash)
	_, err = io.Copy(mw, r)
	if err != nil {
		return "", err
	}
	_, err = tf.Seek(0, 0)
	if err != nil {
		return "", err
	}

	err = v.connect(func(conn *ftp.ServerConn) error {
		if size, err0 := conn.FileSize(filename); err0 == nil && size > 0 {
			return fmt.Errorf("file alredy exist")
		}
		switch dir {
		case ".", "/":
		default:
			if err = conn.MakeDir(dir); err != nil {
				return err
			}
		}
		defer r.Close() //nolint:errcheck

		return conn.Stor(filename, tf)
	})
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (v *FTPProvider) DeleteFile(filename string) error {
	return v.connect(func(conn *ftp.ServerConn) error {
		return conn.Delete(filename)
	})
}
