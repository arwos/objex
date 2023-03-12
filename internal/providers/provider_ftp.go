package providers

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	defer c.Quit()
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
			ctx.Error(http.StatusNotFound, fmt.Errorf("file not found"))
			return nil
		}
		if len(v.conf.Prefix) != 0 {
			ctx.Redirect(v.conf.Prefix + "/" + filename)
			return nil
		}

		resp, err := conn.Retr(filename)
		if err != nil {
			return err
		}
		defer resp.Close()

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

func (v *FTPProvider) SaveFile(filename string, r io.ReadCloser) error {
	dir := filepath.Dir(filename)
	return v.connect(func(conn *ftp.ServerConn) error {
		switch dir {
		case ".", "/":
		default:
			if err := conn.MakeDir(dir); err != nil {
				return err
			}
		}
		defer r.Close()
		return conn.Stor(filename, r)
	})
}

func (v *FTPProvider) DeleteFile(filename string) error {
	return v.connect(func(conn *ftp.ServerConn) error {
		return conn.Delete(filename)
	})
}
