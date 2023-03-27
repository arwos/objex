package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/deweppro/goppy/plugins/web"
)

type (
	tokenContext string
)

const (
	tokenKey = "utk"
)

func TokenDetectMiddleware(cookieName string) web.Middleware {
	detect := []func(r *http.Request) (string, bool){
		func(r *http.Request) (string, bool) {
			value := r.Header.Get("Authorization")
			if len(value) == 0 {
				return "", false
			}
			if strings.Index(value, "Bearer ") == 0 {
				return strings.TrimSpace(value[6:]), true
			}
			return "", false
		},
		func(r *http.Request) (string, bool) {
			if cookie, err := r.Cookie(cookieName); err == nil {
				if len(cookie.Value) > 0 {
					return cookie.Value, true
				}
			}
			return "", false
		},
	}

	return func(call func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
		return func(w http.ResponseWriter, r *http.Request) {
			c := r.Context()
			for _, fn := range detect {
				if token, ok := fn(r); ok {
					c = SetTokenContext(c, token)
					call(w, r.WithContext(c))
					return
				}
			}
			call(w, r)
		}
	}
}

func SetTokenContext(ctx context.Context, value string) context.Context {
	return context.WithValue(ctx, tokenContext(tokenKey), value)
}

func GetTokenContext(c context.Context) (string, bool) {
	value, ok := c.Value(tokenContext(tokenKey)).(string)
	if !ok {
		return "", false
	}
	return value, true
}
