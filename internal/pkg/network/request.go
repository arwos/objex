package network

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type (
	request struct {
		cli *http.Client
		h   http.Header
	}

	Request interface {
		Call(ctx context.Context, r *http.Request, uri string, body io.Reader,
			call func(code int, r io.Reader, codec string) error) error
	}
)

func NewRequest() Request {
	return &request{
		cli: http.DefaultClient,
		h:   make(http.Header),
	}
}

func (v *request) SetHeader(key, value string) {
	v.h.Set(key, value)
}

func (v *request) DeleteHeader(key string) {
	v.h.Del(key)
}

func (v *request) CleanHeaders() {
	for key := range v.h {
		v.h.Del(key)
	}
}

func (v *request) Call(ctx context.Context, r *http.Request, uri string, body io.Reader,
	call func(code int, r io.Reader, codec string) error) error {
	fmt.Println("download", uri, r.Header)
	req, err := http.NewRequestWithContext(ctx, r.Method, uri, body)
	if err != nil {
		return err
	}

	for key := range v.h {
		req.Header.Set(key, v.h.Get(key))
	}

	resp, err := v.cli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint: errcheck

	if err = call(resp.StatusCode, resp.Body, resp.Header.Get("Content-Encoding")); err != nil {
		return err
	}
	return nil
}
