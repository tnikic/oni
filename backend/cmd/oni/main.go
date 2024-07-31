package main

import (
	"github.com/tnikic/oni/server"
	"github.com/tnikic/oni/storage"
)

func main() {
	// Initialize Storage
	var s storage.Provider
	postgres := storage.InitPostgreSQL()
	if postgres != nil {
		s = postgres
	} else {
		sqlite := storage.InitSQLite()
		if sqlite != nil {
			s = sqlite
		} else {
			panic("No storage provider available")
		}
	}

	// Start backend server
	server.Start(s)
}
