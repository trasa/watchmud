package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"html"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/index.html", 302)
}

func someApi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Some Api!! %q", html.EscapeString(r.URL.Path))
}

// Start Up the HttpServer
func Start(port int) {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index)

	// all static goes to /static/ files
	// NOTE: this expects that the program working directory is github.com/trasa/watchmud
	// otherwise, all requests to /static will result in 404!
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static/"))))

	// TODO write an api...
	router.HandleFunc("/api/v1", someApi)

	log.Printf("http listening on port %d", port)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
