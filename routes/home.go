package routes

import "e212/store"

func home(ctx *AppContext) {
	entries, err := store.E212GetAll()
	if err != nil {
		ctx.Error(500, err.Error())
		return
	}
	ctx.Data["need_sorting"] = true
	ctx.Data["nav"] = "home"
	ctx.Data["title"] = "E212 Database"
	ctx.Data["entries"] = entries
	ctx.Data["editentry"] = getCurrentEditItem(ctx)
	ctx.HTML(200, "index")
}
