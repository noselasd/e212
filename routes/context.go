package routes

import (
	"e212/store"

	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

type AppContext struct {
	*macaron.Context
	Flash   *session.Flash
	Session session.Store
}

func SetHeaders() macaron.Handler {
	return func(c *macaron.Context, sess session.Store, f *session.Flash) {
		c.Header().Add("Server", "Macaron Go HTTP")
	}
}

func AppContexter(appVersion string) macaron.Handler {
	return func(c *macaron.Context, f *session.Flash, sess session.Store) {
		ctx := &AppContext{
			Context: c,
			Flash:   f,
			Session: sess,
		}

		if user := sess.Get("user"); user != nil {
			ctx.Data["user"] = user
		}
		c.Data["appVersion"] = appVersion
		c.Map(ctx)

	}
}

func MustBeLoggedIn(ctx *AppContext) {
	var isLoggedIn bool
	user := ctx.Session.Get("user")
	if user != nil {
		if userObj, ok := user.(*store.User); ok {
			isLoggedIn = userObj.ID > 0
		}
	}

	if !isLoggedIn {
		ctx.Error(403, "You are not authorized")
	}

}
