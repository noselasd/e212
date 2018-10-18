package main

import (
	"e212/routes"
	"e212/store"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	//"github.com/go-macaron/macaron"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

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
	routes.InstallRoutes(r)
	runServer(r)
}
