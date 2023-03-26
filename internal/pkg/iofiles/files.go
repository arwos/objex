package iofiles

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	CodecRaw  = ""
	CodecGzip = "gzip"
)

func WriteFile(path string, rc io.Reader, codec string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	dist, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dist.Close() //nolint: errcheck

	switch codec {
	case CodecRaw:
		_, err = io.Copy(dist, rc)
		return err
	case CodecGzip:
		zr, err0 := gzip.NewReader(rc)
		if err0 != nil {
			return err0
		}
		_, err = io.Copy(dist, zr)
		return err
	default:
		return fmt.Errorf("invalid codec")
	}
}
