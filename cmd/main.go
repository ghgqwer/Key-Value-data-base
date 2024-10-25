package main

import (
	"fmt"
	"os"
	"os/signal"
	"project_1/internal/server"
	"project_1/internal/storage/storage"
	"syscall"
)

func main() {
	store, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	// store.Set("first", "maxonP&*r", 5)
	// fmt.Println(store.Get("first"))
	// time.Sleep(5 * time.Second)
	// fmt.Println(store.Get("first"))
	// // time.Sleep(5 * time.Second)
	// // fmt.Println(store.Get("first"))

	// store.Set("first", "maxonP&*r", 0)
	// time.Sleep(3 * time.Second)
	// fmt.Println(store.Get("first"))
	// time.Sleep(5 * time.Second)
	// fmt.Println(store.Get("first"))
	// time.Sleep(10 * time.Second)
	// fmt.Println(store.Get("first"))

	// fmt.Println(store.Lpush("first", []string{"1", "2", "3"}, 10))
	// fmt.Println(store.Rpush("first", []string{"1", "2", "3"}, 0))
	// store.Raddtoset("first", []string{"1", "7"})
	// fmt.Println(store.Check_arr("first"))
	// //store.Raddtoset("first", []string{"1", "7", "123"})
	// store.LSet("first", 0, "19")
	// fmt.Println(store.LGet("first", 0))

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	//gpt :D
	go func() {
		// Блокировка до получения сигнала
		<-c
		fmt.Println("Received shutdown signal")

		// Сохранение данных в файл перед завершением
		if err := store.SaveToJSON("data.json"); err != nil {
			fmt.Printf("Error saving data: %v\n", err)
		} else {
			fmt.Println("Data saved to data.json")
		}

		os.Exit(0) // Завершение программы
	}()

	store.ReadFromJSON("data.json")
	serv := server.New(":8090", &store)
	serv.Start()
}
