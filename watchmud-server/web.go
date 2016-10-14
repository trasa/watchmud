package main

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

func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/static/index.html", 302)
}

func SomeApi(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Some Api!! %q", html.EscapeString(r.URL.Path))
}

func mudsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %s", err)
		return
	}
	client := &Client{conn: c, send: make(chan []byte, 256)}
	go client.writePump()
	client.readPump()
}

func connectHttpServer() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index)

	// all static goes to /static/ files
	// NOTE: this expects that the program working directory is jabmud/jabmud-server
	// otherwise, all requests to /static will result in 404!
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	// TODO write an api...
	router.HandleFunc("/api/v1", SomeApi)

	// websocket handler
	router.HandleFunc("/ws", mudsocket)

	log.Printf("http listening on port %d", port)
	http.Handle("/", router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
