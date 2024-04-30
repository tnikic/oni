package server

import (
	"log"
	"net/http"
)

func StartServer() {
	mux := http.NewServeMux()
	registerHandlers(mux)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func registerHandlers(mux *http.ServeMux) {
	mux.Handle("/", landing())
}
