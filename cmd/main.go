package main

import (
	//"fmt"

	"project_1/internal/server"
	"project_1/internal/storage/storage"
)

func main() {
	store, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	serv := server.New(":8090", &store)
	serv.Start()

}
