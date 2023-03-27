package npm

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/arwos/artifactory/internal/pkg/iofiles"
	"github.com/deweppro/go-sdk/file"
	"github.com/deweppro/go-sdk/log"
	"github.com/deweppro/goppy/plugins/web"
)

func (v *Controller) IndexYarn(c web.Context) {
	hostNpm := v.conf.URISchema() + c.URL().Host + registry
	data := `
yarn config set registryYarn %s

or

YARN_REGISTRY="%s" yarn install

or

.yarnrc:
registryYarn "%s"

`
	c.String(200, data, hostNpm, hostNpm, hostNpm)
}

func (v *Controller) DownloadPackage(c web.Context) {
	path := strings.TrimPrefix(c.URL().Path, registryFiles)
	cacheFile := v.conf.Packages.ProxyCache + path

	if !file.Exist(cacheFile) {
		mux := v.mux.Mutex(path)
		mux.Lock()
		err := v.cli.Call(c.Context(), c.Request(), registryURI+path, nil,
			func(code int, r io.Reader, _ string) error {
				if code != http.StatusOK {
					return fmt.Errorf("status code: %d", code)
				}

				return iofiles.WriteFile(cacheFile, r, iofiles.CodecRaw)
			})
		mux.Unlock()
		if err != nil {
			c.Error(http.StatusBadRequest, err)
			return
		}
	}

	dist, err := os.OpenFile(cacheFile, os.O_RDONLY, 0644)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}
	defer dist.Close() //nolint: errcheck

	c.Response().Header().Set("Content-Type", "application/octet-stream")
	c.Response().WriteHeader(200)
	if _, err = io.Copy(c.Response(), dist); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Errorf("npm files response")
	}
}

func (v *Controller) LoadMetaData(c web.Context) {
	path := strings.TrimPrefix(c.URL().Path, registry)
	metaFile := v.conf.Packages.ProxyCache + path + "/meta.json"

	if !file.Exist(metaFile) {
		mux := v.mux.Mutex(path)
		mux.Lock()
		err := v.cli.Call(c.Context(), c.Request(), registryURI+path, nil,
			func(code int, r io.Reader, codec string) error {
				if code != http.StatusOK {
					return fmt.Errorf("status code: %d", code)
				}

				return iofiles.WriteFile(metaFile, r, codec)
			})
		mux.Unlock()
		if err != nil {
			c.Error(http.StatusBadRequest, err)
			return
		}
	}

	b, err := os.ReadFile(metaFile)
	if err != nil {
		c.Error(500, err)
		return
	}

	hostNpm := v.conf.URISchema() + c.URL().Host + registryFiles
	b = bytes.ReplaceAll(b, []byte(registryURI), []byte(hostNpm))

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	if _, err = c.Response().Write(b); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Errorf("npm files response")
	}
}

func (v *Controller) PublishPackage(c web.Context) {
	publishModel := Publish{}
	err := c.BindJSON(&publishModel)
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}
	fmt.Println(c.Request().Header)
	//fmt.Println(publishModel.Attachments)
	c.String(http.StatusCreated, "ok")
}
