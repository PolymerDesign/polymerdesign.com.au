package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"text/template"
	"time"

	"github.com/gorilla/mux"
)

const hostname = "localhost"
const port = 8080

var debug = false
var cwd string
var templates = template.Must(template.ParseFiles("template/template.html", "template/os.html", "template/home.html"))
var debugFlag = flag.Bool("d", false, "Enable debugging mode.")
var srv *http.Server

func main() {
	var err error
	cwd, err = os.Executable()

	if err != nil {
		log.Fatal(err)
	}

	// Check for debugigng
	flag.Parse()
	if *debugFlag == true {
		fmt.Println("Debugging enabled.")
		debug = true
	}

	createCloseHandler()

	// Create routes
	r := mux.NewRouter()

	r.HandleFunc("/home", PageHandler(HomePageHandler))
	r.HandleFunc("/", RootHandler)

	// Setup server
	srv = &http.Server{
		Handler:      r,
		Addr:         hostname + ":" + strconv.Itoa(port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

// Create close handler for handling CTRL+C events
func createCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// Make sure to clean up properly
	go func() {
		<-c
		fmt.Println("\rCleaning up.")
		srv.Close()
		os.Exit(0)
	}()
}

// RootHandler handles the redirect for the base path
func RootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home", http.StatusPermanentRedirect)
}

// PageHandler is a wrapper for standard page delivery
func PageHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r)
	}
}

// HomePageHandler is the handler for the homepage
func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Debugging bool
	}{
		debug,
	}

	renderTemplate(w, "home", data)
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
	var err error
	var t *template.Template

	if debug {
		t, err = template.ParseFiles("template/template.html", "template/os.html", "template/"+name+".html")
		t.Execute(w, data)
	} else {
		execTemplate(w, data, "template.html", "os.html", name+".html")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func execTemplate(w http.ResponseWriter, data interface{}, names ...string) {
	for _, name := range names {
		err := templates.ExecuteTemplate(w, name, data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
