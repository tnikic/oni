package server

import (
	"log"
	"net/http"

	"github.com/tnikic/oni/controller"
	"github.com/tnikic/oni/storage"
)

var ctr *controller.Controller

func Start(storage storage.Provider) {
	controller.InitDatabase(storage)

	ctr = controller.InitController(storage)

	mux := http.NewServeMux()

	log.Println("Starting server on :80")
	log.Fatal(http.ListenAndServe(":80", mux))
}
