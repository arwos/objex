package npm

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/deweppro/goppy/plugins"
	"github.com/deweppro/goppy/plugins/web"
)

const (
	yarnRegistry = "https://registry.yarnpkg.com"
	npmRegistry  = "https://registry.npmjs.org"
)

var Plugin = plugins.Plugin{
	Config: &Config{},
	Inject: NewController,
}

type Controller struct {
	routes web.RouterPool
	cli    *http.Client
	conf   *Config
}

func NewController(r web.RouterPool, conf *Config) *Controller {
	return &Controller{
		routes: r,
		cli:    http.DefaultClient,
		conf:   conf,
	}
}

func (v *Controller) Up() error {
	route := v.routes.Main()

	route.Get("/yarn/#", v.Yarn)
	route.Get("/npm/#", v.Npm)

	route.Get("/yarn", v.Index)
	route.Get("/npm", v.Index)

	return os.MkdirAll(v.conf.Folder, 0755)
}

func (v *Controller) Down() error {
	return nil
}

func (v *Controller) Index(c web.Context) {
	hostNpm := "http://" + c.URL().Host + "/yarn"
	data := `
yarn config set registry %s

or

YARN_REGISTRY="%s" yarn install

or

.yarnrc:
registry "%s"

`
	c.String(200, data, hostNpm, hostNpm, hostNpm)
}

func (v *Controller) Npm(c web.Context) {
	filename := strings.TrimPrefix(c.URL().Path, "/npm")
	fmt.Println(c.Request().Method, filename)

	req, err := http.NewRequestWithContext(c.Context(), c.Request().Method, npmRegistry+filename, nil)
	if err != nil {
		c.Error(500, err)
		return
	}

	h := c.Request().Header
	for k := range h {
		req.Header.Set(k, h.Get(k))
	}

	resp, err := v.cli.Do(req)
	if err != nil {
		c.Error(500, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.String(resp.StatusCode, "bad request")
		return
	}

	if err = os.MkdirAll(filepath.Dir(v.conf.Folder+filename), 0755); err != nil {
		c.Error(500, err)
		return
	}

	if err = writeBodyToFile(v.conf.Folder+filename, resp.Body, resp.Header.Get("Content-Encoding")); err != nil {
		c.Error(500, err)
		return
	}

	dist, err := os.OpenFile(v.conf.Folder+filename, os.O_RDONLY, 0644)
	if err != nil {
		c.Error(500, err)
		return
	}
	defer dist.Close()

	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().WriteHeader(200)
	if _, err = io.Copy(c.Response(), dist); err != nil {
		fmt.Println(err)
	}
}

func (v *Controller) Yarn(c web.Context) {
	path := strings.TrimPrefix(c.URL().Path, "/yarn")
	fmt.Println(c.Request().Method, path)

	req, err := http.NewRequestWithContext(c.Context(), c.Request().Method, yarnRegistry+path, nil)
	if err != nil {
		c.Error(500, err)
		return
	}

	h := c.Request().Header
	for k := range h {
		req.Header.Set(k, h.Get(k))
	}

	resp, err := v.cli.Do(req)
	if err != nil {
		c.Error(500, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		c.String(resp.StatusCode, "bad request")
		return
	}

	if err = os.MkdirAll(v.conf.Folder+path, 0755); err != nil {
		c.Error(500, err)
		return
	}

	pathMeta := v.conf.Folder + path + "/meta.json"

	if err = writeBodyToFile(pathMeta, resp.Body, resp.Header.Get("Content-Encoding")); err != nil {
		c.Error(500, err)
		return
	}

	b, err := os.ReadFile(pathMeta)
	if err != nil {
		c.Error(500, err)
		return
	}

	hostNpm := "http://" + c.URL().Host + "/npm"
	b = bytes.ReplaceAll(b, []byte(npmRegistry), []byte(hostNpm))

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(resp.StatusCode)
	if _, err = c.Response().Write(b); err != nil {
		fmt.Println(err)
	}
}

func writeBodyToFile(path string, rc io.ReadCloser, codec string) error {
	defer rc.Close()
	dist, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer dist.Close()

	switch codec {
	case "":
		_, err = io.Copy(dist, rc)
		return err
	case "gzip":
		zr, err := gzip.NewReader(rc)
		if err != nil {
			return err
		}
		_, err = io.Copy(dist, zr)
		return err
	default:
		return fmt.Errorf("invalid codec")
	}
}
