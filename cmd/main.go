package main

import (
	"fmt"
	"os"
	"os/signal"
	"project_1/internal/server"
	"project_1/internal/storage/storage"
	"syscall"
	"time"
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

	// store.Set("first", "maxonP&*r", 5)
	// fmt.Println(store.Get("first"))
	// store.Expire("first", 5)
	// time.Sleep(3 * time.Second)
	// fmt.Println(store.Get("first"))

	// store.Lpush("first", []string{"1", "2", "3"}, 10)
	// store.Expire("first", 5)
	// time.Sleep(5 * time.Second)
	// fmt.Println(store.Check_arr("first"))

	closeChan := make(chan struct{})
	go store.GarbageCollection(closeChan, 10*time.Second)

	store.ReadFromJSON("data.json")
	serv := server.New(":8090", &store)
	serv.Start()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	//gpt :D
	go func() {
		// Блокировка до получения сигнала
		<-c
		fmt.Println("Received shutdown signal")
		close(closeChan)
		// Сохранение данных в файл перед завершением
		if err := store.SaveToJSON("data.json"); err != nil {
			fmt.Printf("Error saving data: %v\n", err)
		} else {
			fmt.Println("Data saved to data.json")
		}

		os.Exit(0) // Завершение программы
	}()
}
