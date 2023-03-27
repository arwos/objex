package api

//go:generate static ./../../../web/dist/application ui

import (
	"github.com/deweppro/go-sdk/log"
	"github.com/deweppro/go-static"
	"github.com/deweppro/goppy/plugins/web"
)

var ui static.Reader

func (v *Controller) InjectUIRoutes(route web.RouteCollector) {
	route.Get("/", v.GetUI)
	for _, file := range ui.List() {
		route.Get(file, v.GetUI)
	}
}

func (v *Controller) GetUI(c web.Context) {
	filename := c.URL().Path
	switch filename {
	case "", "/":
		filename = "/index.html"
		break
	}
	if err := ui.ResponseWrite(c.Response(), filename); err != nil {
		log.WithFields(log.Fields{"file": filename, "err": err.Error()}).Errorf("get static file")
	}
}
