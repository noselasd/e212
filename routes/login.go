package routes

import (
	"e212/store"
	"errors"
)

func loginGet(ctx *AppContext) {
	ctx.Data["need_sorting"] = false
	ctx.Data["nav"] = "login"
	ctx.Data["title"] = "Admin Login"

	ctx.HTML(200, "login")
}

func tryLogin(ctx *AppContext) (*store.User, error) {
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

func logout(ctx *AppContext) {
	ctx.Session.Delete("user")
	delete(ctx.Data, "user")
	ctx.Redirect("/")
}

func loginPost(ctx *AppContext) {
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
