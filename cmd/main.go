package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"project_1/internal/filework"
	"project_1/internal/server"
	"project_1/internal/storage/storage"
	"syscall"
	"time"
)

func main() {
	serverPort, ok := os.LookupEnv("BASIC_SERVER_PORT")
	if !ok {
		fmt.Println("not port provided")
		os.Exit(1)
	}

	store, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	if err := filework.ReadFromJSON(store, server.DataJson); err != nil {
		log.Fatal(err)
	}

	closeChan := make(chan struct{})
	go store.GarbageCollection(closeChan, 10*time.Second)
	go store.LoggerSync(closeChan, 10*time.Second)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	//gpt :D
	go func() {
		// Блокировка до получения сигнала
		<-c
		fmt.Println("Received shutdown signal")
		close(closeChan)
		// Сохранение данных в файл перед завершением
		if err := filework.SaveToJSON(store, server.DataJson); err != nil {
			fmt.Printf("Error saving data: %v\n", err)
		} else {
			fmt.Println("Data saved to data.json")
		}

		os.Exit(0) // Завершение программы
	}()

	serv := server.New(":"+serverPort, &store)
	serv.Start()
}
