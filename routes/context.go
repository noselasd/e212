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

func Contexter() macaron.Handler {
	return func(c *macaron.Context, sess session.Store, f *session.Flash) {
		ctx := &AppContext{
			Context: c,
			Flash:   f,
		}

		c.Map(ctx)
	}
}
