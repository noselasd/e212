package routes

import (
	macaron "gopkg.in/macaron.v1"
)

func InstallRoutes(r *macaron.Macaron) {

	r.Get("/", home)
	r.Get("/login", loginGet)
	r.Post("/login", loginPost)
	r.Post("/logout", logout)

	r.Post("/e212update", mustBeLoggedIn, entryUpdate)
	r.Post("/e212add", mustBeLoggedIn, entryAdd)
	r.Post("/e212delete", mustBeLoggedIn, entryDelete)

	r.Group("/e212api.v1/", func() {
		r.Get("/e212", listMCCMNC)
		r.Get("/e212/:mcc", getByMCC)
		r.Get("/e212/:mcc/:mnc", getByMCCMNC)
		r.Delete("/e212/:mcc/:mnc", mustBeLoggedIn, deleteByMCCMNC)
		r.Put("/e212/update", mustBeLoggedIn, updateByMCCMNC)
		r.Post("/e212/create", mustBeLoggedIn, createEntry)
	})

	r.Get("/about", func(ctx *AppContext) {
		otherPagesGet(ctx, "about")
	})
	r.Get("/json-api", func(ctx *AppContext) {
		otherPagesGet(ctx, "json-api")
	})

	r.InternalServerError(e212InternalServerError)
	r.NotFound(e212NotFound)
}
