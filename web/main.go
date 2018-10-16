package main

import (
	"e212/routes"
	"e212/store"
	"errors"

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

func loginGet(ctx *routes.AppContext) {
	ctx.Data["need_sorting"] = false
	ctx.Data["nav"] = "login"
	ctx.Data["title"] = "Admin Login"

	ctx.HTML(200, "login")
}

func tryLogin(ctx *routes.AppContext) (*store.User, error) {
	userName := ctx.QueryTrim("inputUsername")
	passWord := ctx.QueryTrim("inputPassword")

	user, err := store.GetUserByLogin(userName)
	if err == nil {
		if user.CheckPassword(passWord) {
			ctx.Session.Set("user", user)
			ctx.Data["user"] = user
		} else {
			err = errors.New("Password does not match")

		}
	}

	return user, err
}

func logout(ctx *routes.AppContext) {
	ctx.Session.Delete("user")
	delete(ctx.Data, "user")
	ctx.Redirect("/")
}

func loginPost(ctx *routes.AppContext) {
	ctx.Data["need_sorting"] = false
	ctx.Data["nav"] = "login"
	ctx.Data["title"] = "Admin Login"
	_, err := tryLogin(ctx)
	if err != nil {
		ctx.Flash.Error("Unknown user or password", true)
		ctx.HTML(400, "login")
		return
	}
	ctx.Redirect("/")
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
	r.Get("/login", loginGet)
	r.Post("/login", loginPost)
	r.Post("/logout", logout)
	r.Run()
}
