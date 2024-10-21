package main

import (
	"project_1/internal/server"
	"project_1/internal/storage/storage"
)

func main() {
	store, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	store.ReadFromJSON("data.json")
	serv := server.New(":8090", &store)
	serv.Start()
}
