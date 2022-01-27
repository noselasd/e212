package main

import (
	"e212/routes"
	"e212/store"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	//"github.com/go-macaron/macaron"
	"github.com/go-macaron/session"
	"gopkg.in/macaron.v1"
)

var gPort = flag.Int("port", 4000, "Port number to listen on")
var gUseTLS = flag.Bool("usetls", false, "Use TLS(HTTPS) instead of plain HTTP")
var gTLSCert = flag.String("tlscert", "tls.cert", "Path to TLS certificate file")
var gTLSKey = flag.String("tlskey", "tls.key", "Path to TLS key file")
var gTraceSQL = flag.Bool("logsql", false, "Log SQL statements")
var gVersion = "DEVELOPMENT"

func runServer(r *macaron.Macaron) {

	srv := http.Server{
		Addr:         fmt.Sprintf(":%d", *gPort),
		Handler:      r,
		ReadTimeout:  45 * time.Second,
		WriteTimeout: 45 * time.Second,
	}

	log.Printf("v. %s listening on %s\n", gVersion, srv.Addr)

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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s v.%s\n", os.Args[0], gVersion)
	}

	flag.Parse()
	err := store.Init("mccmnc.db", *gTraceSQL)
	if err != nil {
		panic(err)
	}
	m := macaron.New()

	m.Use(macaron.Logger()) 		//standard logger
	m.Use(macaron.Recovery()) 		//standard recovery if code panics
	m.Use(macaron.Static("public")) //serve static files from public/ folder mapped to /public
	m.Use(macaron.Renderer()) 		//standard html template renderer
	m.Use(session.Sessioner()) 		//standard http session support
	m.Use(routes.SetHeaders())  	//our custom http headers
	m.Use(routes.AppContexter(gVersion))//aplication context to pass to route handlers
	routes.InstallRoutes(m) 		//Our routes
	runServer(m)
}
