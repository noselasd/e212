package main

import (
	"e212/routes"
	"e212/store"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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

func errRedirect(ctx *routes.AppContext, location string, errMsg string) {
	ctx.Flash.Error(errMsg)
	ctx.Redirect(location)
}

func entryDelete(ctx *routes.AppContext) {
	id := ctx.QueryTrim("inputID")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		errRedirect(ctx, "/", "Update failed: "+err.Error())
		return
	}

	err = store.E212DeleteById(idInt)
	if err != nil {
		errRedirect(ctx, "/", "Delete failed: "+err.Error())
		return
	}

	ctx.Redirect("/")
}

func entryUpdate(ctx *routes.AppContext) {
	handleAddEdit(ctx, false)
}

func entryAdd(ctx *routes.AppContext) {
	handleAddEdit(ctx, true)
}

func handleAddEdit(ctx *routes.AppContext, isNew bool) {
	id := ctx.QueryTrim("inputID")
	country := ctx.QueryTrim("inputCountry")
	operator := ctx.QueryTrim("inputOperator")
	mcc := ctx.QueryTrim("inputMCC")
	mnc := ctx.QueryTrim("inputMNC")

	entry := store.NewE212Entry(mcc, mnc, country, operator)
	err := entry.Validate()
	if err != nil {
		errRedirect(ctx, "/", "Update failed: "+err.Error())
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		errRedirect(ctx, "/", "Update failed: "+err.Error())
		return
	}
	entry.ID = idInt
	if isNew {
		err = store.E212Add(entry)
	} else {
		err = store.E212Update(entry)

	}
	if err != nil {
		errRedirect(ctx, "/", "Update failed: "+err.Error())
		return
	}

	//ok
	ctx.Redirect("/")
}

var gPort = flag.Int("port", 4000, "port number to listen on")
var gUseTLS = flag.Bool("usetls", false, "Use TLS(HTTPS) intead of plain HTTP")
var gTLSCert = flag.String("tlscert", "tls.cert", "Path to TLS certificate file")
var gTLSKey = flag.String("tlskey", "tls.key", "Path to TLS key file")

func runServer(r *macaron.Macaron) {

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", *gPort),
		Handler:      r,
		ReadTimeout:  45 * time.Second,
		WriteTimeout: 45 * time.Second,
	}

	log.Printf("listening on %s\n", srv.Addr)

	var err error
	if *gUseTLS {
		err = srv.ListenAndServeTLS(*gTLSCert, *gTLSKey)
	} else {
		err = srv.ListenAndServe()
	}
	if err != nil {
		log.Fatalln("ListenAndServe: ", err)
	}

}

func main() {
	flag.Parse()
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

	r.Post("/e212update", routes.MustBeLoggedIn, entryUpdate)
	r.Post("/e212add", routes.MustBeLoggedIn, entryAdd)
	r.Post("/e212delete", routes.MustBeLoggedIn, entryDelete)

	runServer(r)
}
