package controller

//go:generate easyjson

import (
	"net/http"

	"github.com/deweppro/goppy/plugins/web"
)

//easyjson:json
type NewUserModel struct {
	Login  string `json:"login"`
	Passwd string `json:"passwd"`
}

func (v *Controller) CreateUser(c web.Context) {
	model := NewUserModel{}
	if err := c.BindJSON(&model); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	err := v.users.CreateUser(c.Context(), model.Login, model.Passwd)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusOK, "ok")
}

//easyjson:json
type NewGroupModel struct {
	Name string `json:"name"`
}

func (v *Controller) CreateGroup(c web.Context) {
	model := NewGroupModel{}
	if err := c.BindJSON(&model); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	err := v.users.CreateGroup(c.Context(), model.Name)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusOK, "ok")
}

type (
	//easyjson:json
	ListGroupsModel []ListGroupModel
	//easyjson:json
	ListGroupModel struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
)

func (v *Controller) ListGroup(c web.Context) {
	model := make(ListGroupsModel, 0)

	list, err := v.users.ListGroup(c.Context())
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	for id, name := range list {
		model = append(model, ListGroupModel{
			ID:   id,
			Name: name,
		})
	}

	c.JSON(http.StatusOK, model)
}

//easyjson:json
type UserAddGroupModel struct {
	Login string  `json:"login"`
	IDs   []int64 `json:"ids"`
}

func (v *Controller) AddUserGroup(c web.Context) {
	model := UserAddGroupModel{}
	if err := c.BindJSON(&model); err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	err := v.users.AppendUserToGroups(c.Context(), model.Login, model.IDs...)
	if err != nil {
		c.Error(http.StatusBadRequest, err)
		return
	}

	c.String(http.StatusOK, "ok")
}
