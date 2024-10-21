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
