package routes

import macaron "gopkg.in/macaron.v1"

func InstallRoutes(r *macaron.Macaron) {
	r.Group("/e212api.v1/", func() {
		r.Get("/e212", ListMCCMNC)
		r.Get("/e212/:mcc", GetByMCC)
		r.Get("/e212/:mcc/:mnc", GetByMCCMNC)
	})

	r.Get("/", home)
	r.Get("/login", loginGet)
	r.Post("/login", loginPost)
	r.Post("/logout", logout)

	r.Post("/e212update", MustBeLoggedIn, entryUpdate)
	r.Post("/e212add", MustBeLoggedIn, entryAdd)
	r.Post("/e212delete", MustBeLoggedIn, entryDelete)

	r.Group("/e212api.v1/", func() {
		r.Get("/e212", ListMCCMNC)
		r.Get("/e212/:mcc", GetByMCC)
		r.Get("/e212/:mcc/:mnc", GetByMCCMNC)
	})

	r.Get("/", home)
	r.Get("/login", loginGet)
	r.Post("/login", loginPost)
	r.Post("/logout", logout)

	r.Post("/e212update", MustBeLoggedIn, entryUpdate)
	r.Post("/e212add", MustBeLoggedIn, entryAdd)
	r.Post("/e212delete", MustBeLoggedIn, entryDelete)
}