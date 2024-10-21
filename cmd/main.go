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

	// store.Rpush("first", []string{"1", "2", "3"})
	// fmt.Println(store.Check_arr("first"))
	// store.Lpop("first", []int{1})
	// fmt.Println(store.Check_arr("first"))

	store.ReadFromJSON("data.json")
	serv := server.New(":8090", &store)
	serv.Start()
}
