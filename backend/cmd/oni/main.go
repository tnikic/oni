package main

import (
	"github.com/tnikic/oni/server"
	"github.com/tnikic/oni/storage"
)

func main() {
	// Initialize Storage
	storage := storage.InitSQLite()

	// Start backend server
	server.Start(storage)
}
