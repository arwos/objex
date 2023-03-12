package controller

import (
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

	prov, err := v.providers.GetByCode(store.Code)
	if err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}

	if err = prov.SaveFile(filename, c.Request().Body); err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}

	if err = v.files.AddFile(c.Context(), store.ID, filename, "", c.URL().Query()); err != nil {
		c.Error(http.StatusInternalServerError, err)
		return
	}

	c.String(http.StatusOK, "ok")
}
