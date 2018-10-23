package routes

import (
	"e212/store"
	"errors"
	"time"

	macaron "gopkg.in/macaron.v1"
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

	user, ok := ctx.Session.Get("user").(*store.User)

	ctx.Session.Delete("user")
	delete(ctx.Data, "user")

	if ok {
		ctx.Logger.Printf("%s: User %s logged out\n", time.Now().Format(macaron.LogTimeFormat), user.LoginName)
	}
	ctx.Redirect("/")
}

func loginPost(ctx *AppContext) {
	ctx.Data["need_sorting"] = false
	ctx.Data["nav"] = "login"
	ctx.Data["title"] = "Admin Login"
	user, err := tryLogin(ctx)
	if err != nil {
		ctx.Flash.Error("Unknown user or password", true)
		ctx.Logger.Printf("%s: User %s failed login\n", time.Now().Format(macaron.LogTimeFormat), ctx.QueryTrim("inputUsername"))
		ctx.HTML(400, "login")
		return
	}
	ctx.Logger.Printf("%s: User %s logged in\n", time.Now().Format(macaron.LogTimeFormat), user.LoginName)
	ctx.Redirect("/")
}
