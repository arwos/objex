package controller

import (
	"fmt"
	"net/http"

	"github.com/deweppro/goppy/plugins/web"
)

func (v *Controller) UploadFile(c web.Context) {
	login, passwd, ok := c.Request().BasicAuth()
	if !ok {
		c.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(c.Response(), "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !v.users.ValidateUserPasswd(c.Context(), login, passwd) {
		c.String(http.StatusForbidden, "Unauthorized")
		return
	}

	storeName, err := c.Param("storage").String()
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	filename := c.URL().Path[len(storeName)+2:]

	store, err := v.store.Get(c.Context(), storeName)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	if !v.users.HasUserInGroup(c.Context(), login, store.GetGroups()...) {
		c.String(http.StatusForbidden, "Access denied")
		return
	}

	prov, err := v.providers.GetByCode(store.Code)
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}
	hash, err := prov.SaveFile(filename, c.Request().Body)
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}

	if err = v.files.AddFile(c.Context(), store.ID, filename, hash, c.URL().Query()); err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "ok")
}

func (v *Controller) DownloadFile(c web.Context) {
	login, passwd, ok := c.Request().BasicAuth()
	if !ok {
		c.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(c.Response(), "Unauthorized", http.StatusUnauthorized)
		return
	}

	if !v.users.ValidateUserPasswd(c.Context(), login, passwd) {
		c.String(http.StatusForbidden, "Unauthorized")
		return
	}

	storeName, err := c.Param("storage").String()
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	filename := c.URL().Path[len(storeName)+2:]

	store, err := v.store.Get(c.Context(), storeName)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	if !v.users.HasUserInGroup(c.Context(), login, store.GetGroups()...) {
		c.String(http.StatusForbidden, "Access denied")
		return
	}

	if !v.files.HasFile(c.Context(), store.ID, filename) {
		c.Error(http.StatusNotFound, fmt.Errorf("file not found"))
		return
	}

	prov, err := v.providers.GetByCode(store.Code)
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}

	prov.GetFile(filename, c)
}

func (v *Controller) SearchByProps(c web.Context) {

}
