package routes

import (
	"e212/store"
)

type jsonErr struct {
	Error string `json:"errormessage"`
}

func jsonError(httpStatus int, err error, ctx *AppContext) {
	e := jsonErr{Error: err.Error()}
	ctx.Logger.Printf("%s %s failed: %v\n ", ctx.Req.Method, ctx.Req.URL, e)
	ctx.JSON(httpStatus, &e)
}

func ListMCCMNC(ctx *AppContext) {
	e, err := store.E212GetAll()
	if err != nil {
		jsonError(500, err, ctx)
	} else {
		ctx.JSON(200, e)
	}
}

func GetByMCC(ctx *AppContext) {

	mcc := ctx.Params("mcc")

	e, err := store.E212GetByMcc(mcc)
	if err == nil {
		ctx.JSON(200, e)
	} else {
		jsonError(404, err, ctx)
	}
}

func GetByMCCMNC(ctx *AppContext) {
	mcc := ctx.Params("mcc")
	mnc := ctx.Params("mnc")

	e, err := store.E212GetByMccMnc(&store.MccMnc{Mcc: mcc, Mnc: mnc})
	if err == nil {
		ctx.JSON(200, e)
	} else {
		jsonError(404, err, ctx)
	}
}
