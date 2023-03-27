package api

import (
	"net/http"

	"github.com/arwos/artifactory/internal/pkg/middlewares"
	"github.com/deweppro/go-sdk/log"
	"github.com/deweppro/goppy/plugins/web"
)

func (v *Controller) InjectAuthRoutes(route web.RouteCollector) {
	route.Post("/login", v.AuthLogin)
	route.Delete("/logout", v.AuthLogout)
	route.Get("/check", v.AuthCheck)
}

func (v *Controller) AuthLogin(c web.Context) {
	if token, ok := middlewares.GetTokenContext(c.Context()); ok {
		if _, err := v.users.GetUserByToken(c.Context(), token); err != nil {
			c.JSON(http.StatusOK, EMPTY)
			return
		}
	}

	model := LoginRequest{}
	if err := c.BindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, EMPTY)
		return
	}
	if len(model.Login) == 0 || len(model.Password) == 0 {
		c.JSON(http.StatusForbidden, EMPTY)
		return
	}

	if !v.users.ValidateUserPasswd(c.Context(), model.Login, model.Password) {
		c.JSON(http.StatusForbidden, EMPTY)
		return
	}
	token, err := v.users.CreateToken(c.Context(), model.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, EMPTY)
		return
	}
	cookie := &http.Cookie{
		Name:  v.conf.Settings.CookieName,
		Value: token,
		//Domain:   c.URL().Host,
		Secure:   v.conf.Settings.UsedHTTPS,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	}
	c.Cookie().Set(cookie)
	c.JSON(http.StatusOK, EMPTY)
}

func (v *Controller) AuthLogout(c web.Context) {
	token, ok := middlewares.GetTokenContext(c.Context())
	if !ok {
		c.JSON(http.StatusOK, EMPTY)
		return
	}
	u, err := v.users.GetUserByToken(c.Context(), token)
	if err != nil {
		log.WithFields(log.Fields{"token": err.Error()}).Errorf("logout")
		c.JSON(http.StatusOK, EMPTY)
		return
	}
	if err = v.users.DeleteToken(c.Context(), u.Login, token); err != nil {
		log.WithFields(log.Fields{"token": err.Error(), "uid": u.ID}).Errorf("logout")
	}
	c.JSON(http.StatusOK, EMPTY)
}

func (v *Controller) AuthCheck(c web.Context) {
	if token, ok := middlewares.GetTokenContext(c.Context()); ok {
		if _, err := v.users.GetUserByToken(c.Context(), token); err == nil {
			c.JSON(http.StatusOK, EMPTY)
			return
		}
	}
	c.JSON(http.StatusUnauthorized, EMPTY)
}
