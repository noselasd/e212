package routes

import (
	"e212/store"

	"gopkg.in/macaron.v1"
)

type jsonErr struct {
	Error string `json:"errormessage"`
}

func jsonError(httpStatus int, err error, ctx *macaron.Context) {
	e := jsonErr{Error: err.Error()}

	ctx.JSON(httpStatus, &e)
}

func ListMCCMNC(ctx *macaron.Context) {
	e, err := store.E212GetAll()
	if err != nil {
		jsonError(500, err, ctx)
	} else {
		ctx.JSON(200, e)
	}
}

func GetByMCC(ctx *macaron.Context) {

	mcc := ctx.Params("mcc")

	e, err := store.E212GetByMcc(mcc)
	if err == nil {
		ctx.JSON(200, e)
	} else {
		jsonError(404, err, ctx)
	}
}

func GetByMCCMNC(ctx *macaron.Context) {
	mcc := ctx.Params("mcc")
	mnc := ctx.Params("mnc")

	e, err := store.E212GetByMccMnc(&store.MccMnc{Mcc: mcc, Mnc: mnc})
	if err == nil {
		ctx.JSON(200, e)
	} else {
		jsonError(404, err, ctx)
	}
}
