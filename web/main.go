package main

import (
	"e212/routes"
	"e212/store"

	//"github.com/go-macaron/macaron"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

func home(ctx *routes.AppContext) {
	entries, err := store.E212GetAll()
	if err != nil {
		ctx.Error(500, err.Error())
	}
	ctx.Data["need_sorting"] = true
	ctx.Data["nav"] = "home"
	ctx.Data["title"] = "E212 Database"
	ctx.Data["entries"] = entries
	ctx.HTML(200, "index")
}

func main() {
	err := store.Init("mccmnc.db")
	if err != nil {
		panic(err)
	}
	r := macaron.Classic()
	r.Use(macaron.Renderer())
	r.Use(session.Sessioner())
	r.Use(routes.SetHeaders())
	r.Use(routes.AppContexter())

	r.Group("/e212api.v1/", func() {
		r.Get("/e212", routes.ListMCCMNC)
		r.Get("/e212/:mcc", routes.GetByMCC)
		r.Get("/e212/:mcc/:mnc", routes.GetByMCCMNC)
	})

	r.Get("/", home)
	r.Run()
}
