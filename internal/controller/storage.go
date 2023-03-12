package controller

import (
	"net/http"

	"github.com/deweppro/goppy/plugins/web"
)

//easyjson:json
type NewStorageModel struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Lifetime int64  `json:"lifetime"`
}

func (v *Controller) CreateStorage(c web.Context) {
	model := NewStorageModel{}
	if err := c.BindJSON(&model); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	err := v.store.CreateStore(c.Context(), model.Name, model.Code, model.Lifetime)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusOK, "ok")
}

//easyjson:json
type AddStorageGroupModel struct {
	Name string  `json:"name"`
	IDs  []int64 `json:"ids"`
}

func (v *Controller) AddStorageGroup(c web.Context) {
	model := AddStorageGroupModel{}
	if err := c.BindJSON(&model); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	err := v.store.AppendStorageToGroups(c.Context(), model.Name, model.IDs...)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusOK, "ok")
}
