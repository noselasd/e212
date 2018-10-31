package routes

import (
	"e212/store"
	"strconv"
)

func entryDelete(ctx *AppContext) {
	ctx.Data["editentry"] = getCurrentEditItem(ctx)

	id := ctx.QueryTrim("inputID")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errRedirect(ctx, "/", "Delete failed: "+err.Error())
		return
	}

	err = store.E212DeleteById(idInt)
	if err != nil {
		errRedirect(ctx, "/", "Delete failed: "+err.Error())
		return
	}

	ctx.Redirect("/")
}

func entryUpdate(ctx *AppContext) {
	handleAddEdit(ctx, false)
}

func entryAdd(ctx *AppContext) {
	handleAddEdit(ctx, true)
}

func handleAddEdit(ctx *AppContext, isNew bool) {
	id := ctx.QueryTrim("inputID")
	country := ctx.QueryTrim("inputCountry")
	operator := ctx.QueryTrim("inputOperator")
	mcc := ctx.QueryTrim("inputMCC")
	mnc := ctx.QueryTrim("inputMNC")

	entry := store.NewE212Entry(mcc, mnc, country, operator)
	setCurrentEditItem(ctx, entry)

	err := entry.Validate()
	if err != nil {
		errRedirect(ctx, "/", "Operation failed: "+err.Error())
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		errRedirect(ctx, "/", "Operation failed: "+err.Error())
		return
	}
	entry.ID = idInt
	if isNew {
		err = store.E212Add(entry)
	} else {
		err = store.E212Update(entry)

	}
	if err != nil {
		errRedirect(ctx, "/", "Operation failed: "+err.Error())
		return
	}

	//ok
	ctx.Redirect("/")
}
