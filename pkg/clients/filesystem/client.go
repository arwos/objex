package filesystem

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"

	"go.arwos.org/objex/pkg/clients"
	"go.osspkg.com/ioutils/fs"
)

type Client struct {
	folder string

	volumeSize     int64
	volumeBusySize int64
}

func (c *Client) Close() (err error) {
	return nil
}

func (c *Client) calcVolumeSize() (err error) {
	var size int64
	err = filepath.Walk(c.folder, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	if err != nil {
		return err
	}
	atomic.StoreInt64(&c.volumeBusySize, size)
	return nil
}

func (c *Client) Connect(address string, opt ...any) error {
	if !fs.FileExist(address) {
		err := os.MkdirAll(address, 0777)
		if err != nil {
			return fmt.Errorf("filesystem connect: %w", err)
		}
	}

	opts, err := clients.NewOptions(opt)
	if err != nil {
		return fmt.Errorf("filesystem connect: %w", err)
	}

	c.folder = address
	c.volumeSize = opts.IntValue("volume_size", 0)

	if err = c.calcVolumeSize(); err != nil {
		return fmt.Errorf("filesystem connect: %w", err)
	}

	return nil
}

func (c *Client) fullPath(filename string) string {
	return filepath.Join(c.folder, filename)
}

func (c *Client) Delete(filename string) (err error) {
	path := c.fullPath(filename)

	fi, err := os.Stat(path)
	if err != nil {
		return err
	}

	if err = os.RemoveAll(path); err != nil {
		return err
	}

	atomic.AddInt64(&c.volumeBusySize, -fi.Size())

	return nil
}

func (c *Client) Exist(filename string) (ok bool) {
	return fs.FileExist(c.fullPath(filename))
}

func (c *Client) Free() (n int64, err error) {
	n = c.volumeSize - atomic.LoadInt64(&c.volumeBusySize)
	return n, nil
}

func (c *Client) Load(filename string) (r io.ReadCloser, err error) {
	fo, err := os.Open(c.fullPath(filename))
	if err != nil {
		return nil, err
	}
	return fo, nil
}

func (c *Client) Save(filename string, r io.Reader) (err error) {
	fc, err := os.Create(c.fullPath(filename))
	if err != nil {
		return err
	}

	size, err := io.Copy(fc, r)
	if err != nil {
		return err
	}

	atomic.AddInt64(&c.volumeBusySize, size)

	return nil
}

func (c *Client) Size(filename string) (n int64, err error) {
	fi, err := os.Stat(c.fullPath(filename))
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}
