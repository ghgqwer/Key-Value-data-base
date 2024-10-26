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

	closeChan := make(chan struct{})
	go store.GarbageCollection(closeChan, 10*time.Second)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	//gpt :D
	go func() {
		// Блокировка до получения сигнала
		<-c
		fmt.Println("Received shutdown signal")
		close(closeChan)
		// Сохранение данных в файл перед завершением
		if err := store.SaveToJSON(server.DataJson); err != nil {
			fmt.Printf("Error saving data: %v\n", err)
		} else {
			fmt.Println("Data saved to data.json")
		}

		os.Exit(0) // Завершение программы
	}()

	store.ReadFromJSON(server.DataJson)
	serv := server.New(":8090", &store)
	serv.Start()
}
