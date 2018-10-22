package routes

import (
	"log"

	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

type AppContext struct {
	*macaron.Context
	Flash   *session.Flash
	Session session.Store
	Logger  *log.Logger
}

func SetHeaders() macaron.Handler {
	return func(c *macaron.Context, sess session.Store, f *session.Flash) {
		c.Header().Add("Server", "Macaron Go HTTP")
	}
}

func AppContexter(appVersion string) macaron.Handler {
	return func(c *macaron.Context, f *session.Flash, sess session.Store, log *log.Logger) {
		ctx := &AppContext{
			Context: c,
			Flash:   f,
			Session: sess,
			Logger:  log,
		}

		if user := sess.Get("user"); user != nil {
			ctx.Data["user"] = user
		}
		c.Data["appVersion"] = appVersion
		c.Map(ctx)

	}
}
