package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"html"
	"log"
	"net/http"
)

const port = 8888

var upgrader = websocket.Upgrader{} // default options

func index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/index.html", 302)
}

func someApi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Some Api!! %q", html.EscapeString(r.URL.Path))
}

func mudsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %s", err)
		return
	}
	// new client object, add to clients and
	// set up reader/writer
	client := newWebClient(c)
	go client.writePump()
	client.readPump()
}

// Start Up the HttpServer
func Start() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", index)

	// all static goes to /static/ files
	// NOTE: this expects that the program working directory is github.com/trasa/watchmud
	// otherwise, all requests to /static will result in 404!
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./server/static/"))))

	// TODO write an api...
	router.HandleFunc("/api/v1", someApi)

	// websocket handler
	router.HandleFunc("/ws", mudsocket)

	log.Printf("http listening on port %d", port)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}