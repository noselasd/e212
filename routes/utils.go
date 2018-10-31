package routes

import (
	"e212/store"
	"errors"
	"net/http"
	"strings"
	"time"

	macaron "gopkg.in/macaron.v1"
)

func acceptsJson(ctx *AppContext) bool {
	accept := ctx.Req.Header.Get("Accept")
	return strings.Index(accept, "application/json") != -1
}

func e212InternalServerError(ctx *AppContext, err error) {

	if acceptsJson(ctx) {
		jsonError(500, err, ctx)
	} else {
		http.Error(ctx.Resp, err.Error(), 500)
	}
}

func e212NotFound(ctx *AppContext, req *http.Request) {

	if acceptsJson(ctx) {
		jsonError(404, errors.New("Not Found"), ctx)
	} else {
		http.NotFound(ctx.Resp, req)
	}
}

func mustBeLoggedIn(ctx *AppContext) {
	var isLoggedIn bool
	user := ctx.Session.Get("user")
	if user != nil {
		if userObj, ok := user.(*store.User); ok {
			isLoggedIn = userObj.ID > 0
		}
	}

	if !isLoggedIn {
		err := errors.New("You are not authorized")
		if acceptsJson(ctx) {
			jsonError(500, err, ctx)
		} else {
			ctx.Error(403, err.Error())
		}
	}
}

func getCurrentEditItem(ctx *AppContext) *store.E212Entry {
	e := ctx.Session.Get("editentry")
	if e != nil {
		if eObj, ok := e.(*store.E212Entry); ok {
			return eObj
		}
	}

	return &store.E212Entry{}
}

func setCurrentEditItem(ctx *AppContext, e *store.E212Entry) {
	ctx.Session.Set("editentry", e)
}

func errRedirect(ctx *AppContext, location string, errMsg string) {
	ctx.Flash.Error(errMsg)
	ctx.Header().Set("Warning", errMsg)
	ctx.Header().Set("Status", "400 request error")

	ctx.Logger.Printf("%s: %s %s failed: %s\n ", time.Now().Format(macaron.LogTimeFormat), ctx.Req.Method, ctx.Req.URL, errMsg)

	ctx.Redirect(location)
}
