package routes

import (
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

func AppContexter() macaron.Handler {
	return func(c *macaron.Context, f *session.Flash, sess session.Store) {
		ctx := &AppContext{
			Context: c,
			Flash:   f,
			Session: sess,
		}

		if user := sess.Get("user"); user != nil {
			ctx.Data["user"] = user
		}

		c.Map(ctx)

	}
}
