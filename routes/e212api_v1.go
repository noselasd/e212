package routes

import (
	"e212/store"
	"encoding/json"
	"errors"
	"time"

	macaron "gopkg.in/macaron.v1"
)

type jsonErr struct {
	Error string `json:"errormessage"`
}

func jsonError(httpStatus int, err error, ctx *AppContext) {
	e := jsonErr{Error: err.Error()}
	ctx.Logger.Printf("%s: %s %s failed: %v\n ", time.Now().Format(macaron.LogTimeFormat), ctx.Req.Method, ctx.Req.URL, e)
	ctx.JSON(httpStatus, &e)
}

func listMCCMNC(ctx *AppContext) {
	e, err := store.E212GetAll()
	if err != nil {
		jsonError(500, err, ctx)
	} else {
		ctx.JSON(200, e)
	}
}

func getByMCC(ctx *AppContext) {

	mcc := ctx.Params("mcc")

	e, err := store.E212GetByMcc(mcc)
	if err == nil {
		ctx.JSON(200, e)
	} else {
		jsonError(404, err, ctx)
	}
}

func getByMCCMNC(ctx *AppContext) {
	mcc := ctx.Params("mcc")
	mnc := ctx.Params("mnc")

	e, err := store.E212GetByMccMnc(&store.MccMnc{Mcc: mcc, Mnc: mnc})
	if err == nil {
		ctx.JSON(200, e)
	} else {
		jsonError(404, err, ctx)
	}
}

func deleteByMCCMNC(ctx *AppContext) {
	mcc := ctx.Params("mcc")
	mnc := ctx.Params("mnc")

	e, err := store.E212GetByMccMnc(&store.MccMnc{Mcc: mcc, Mnc: mnc})
	if err != nil {
		jsonError(404, err, ctx)
		return
	}

	err = store.E212DeleteById(e.ID)
	if err == store.ErrEntryMissing {
		jsonError(404, err, ctx)
		return
	} else if err != nil {
		jsonError(500, err, ctx)
		return
	}

	ctx.Status(204)
}

func deleteById(ctx *AppContext) {
	id := ctx.ParamsInt("id")

	err := store.E212DeleteById(id)
	if err == store.ErrEntryMissing {
		jsonError(404, err, ctx)
		return
	} else if err != nil {
		jsonError(500, err, ctx)
		return
	}

	ctx.Status(204)
}

func readJsonEntry(ctx *AppContext) (*store.E212Entry, error) {
	bodyReader := ctx.Req.Body().ReadCloser()
	defer bodyReader.Close()

	decoder := json.NewDecoder(bodyReader)
	var entry store.E212Entry

	err := decoder.Decode(&entry)

	return &entry, err
}

func updateByMCCMNC(ctx *AppContext) {

	entry, err := readJsonEntry(ctx)
	if err != nil {
		jsonError(400, err, ctx)
		return
	}

	if err = entry.Validate(); err != nil {
		jsonError(400, errors.New("Validation error"), ctx)
		return
	}

	e, err := store.E212GetByMccMnc(&entry.Code)
	if err != nil {
		jsonError(404, err, ctx)
		return
	}
	entry.ID = e.ID

	err = store.E212Update(entry)
	if err != nil {
		jsonError(500, err, ctx)
		return
	}

	ctx.Status(204)
}

func createEntry(ctx *AppContext) {

	entry, err := readJsonEntry(ctx)
	if err != nil {
		jsonError(400, err, ctx)
		return
	}

	if err = entry.Validate(); err != nil {
		jsonError(400, errors.New("Validation error"), ctx)
		return
	}
	entry.ID = 0

	err = store.E212Add(entry)
	if err != nil {
		jsonError(500, err, ctx)
		return
	}

	ctx.Status(204)
}
