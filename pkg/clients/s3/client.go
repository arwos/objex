package s3

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.arwos.org/objex/pkg/clients"
	"go.osspkg.com/static"
)

type Client struct {
	cli        *minio.Client
	bucketName string
	location   string

	volumeSize     int64
	volumeBusySize int64
}

func (c *Client) Close() (err error) {
	return nil
}

func (c *Client) calcVolumeSize() error {
	var size int64

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	objectCh := c.cli.ListObjects(ctx, c.bucketName, minio.ListObjectsOptions{Recursive: true})
	for object := range objectCh {
		if object.Err != nil {
			return object.Err
		}
		size += object.Size
	}
	atomic.StoreInt64(&c.volumeBusySize, size)
	return nil
}

func (c *Client) Connect(address string, opt ...any) error {
	opts, err := clients.NewOptions(opt)
	if err != nil {
		return fmt.Errorf("s3 connect: %w", err)
	}

	accessKeyID := opts.StringValue("access_key_id", "")
	secretAccessKey := opts.StringValue("secret_access_key", "")
	useSSL := opts.BoolValue("secret_access_key", false)
	c.bucketName = opts.StringValue("bucket_name", "")
	c.location = opts.StringValue("location", "")
	c.volumeSize = opts.IntValue("volume_size", 0)

	c.cli, err = minio.New(address, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return fmt.Errorf("s3 connect: %w", err)
	}

	if err = c.calcVolumeSize(); err != nil {
		return fmt.Errorf("s3 connect: %w", err)
	}

	return nil
}

func (c *Client) Delete(filename string) (err error) {
	size, err := c.Size(filename)
	if err != nil {
		return err
	}
	err = c.cli.RemoveObject(context.TODO(), c.bucketName, filename, minio.RemoveObjectOptions{ForceDelete: true, GovernanceBypass: true})
	if err != nil {
		return err
	}
	atomic.AddInt64(&c.volumeBusySize, -size)
	return nil
}

func (c *Client) Exist(filename string) (ok bool) {
	n, err := c.Size(filename)
	return err != nil || n == 0
}

func (c *Client) Free() (n int64, err error) {
	n = c.volumeSize - atomic.LoadInt64(&c.volumeBusySize)
	return n, nil
}

func (c *Client) Load(filename string) (io.ReadCloser, error) {
	r, err := c.cli.GetObject(context.TODO(), c.bucketName, filename, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	w, err := os.CreateTemp("/tmp", "objex-*.bin")
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(w, r); err != nil {
		w.Close()
		os.RemoveAll(w.Name())
		return nil, err
	}
	if _, err = w.Seek(0, 0); err != nil {
		w.Close()
		os.RemoveAll(w.Name())
		return nil, err
	}
	return w, nil
}

func (c *Client) Save(filename string, r io.Reader) error {
	contentType := static.DetectContentType(filename, nil)
	info, err := c.cli.PutObject(
		context.TODO(),
		c.bucketName,
		filename,
		r,
		-1,
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		return err
	}
	atomic.AddInt64(&c.volumeBusySize, info.Size)
	return nil
}

func (c *Client) Size(filename string) (n int64, err error) {
	info, err := c.cli.StatObject(context.TODO(), c.bucketName, filename, minio.StatObjectOptions{})
	if err != nil {
		return 0, nil
	}
	return info.Size, nil
}
