package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"project_1/internal/filework"
	"project_1/internal/server"
	"project_1/internal/storage/storage"
	"strconv"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

const (
	queryCreateStatesTable = `CREATE TABLE IF NOT EXISTS core (
		version bigserial PRIMARY KEY,
		timestamp bigint NOT NULL,
		payload JSONB NOT NULL
	)`

	quertCreateState = `INSERT INTO core (timestamp, payload) VALUES ($1, $2)`
)

type Core struct {
	Version   int    `json: version`
	Timestamp int    `json: timestamp`
	Payload   string `json: payload`
}

func saveToDB(db *sql.DB, storage storage.Storage, timestamp int64, closeChan <-chan struct{}, n time.Duration) {
	for {
		select {
		case <-closeChan:
			return
		case <-time.After(n):
			payload, _ := json.Marshal(storage)
			db.Exec(quertCreateState, timestamp, payload)

			db.Exec(`DELETE FROM core
							  WHERE version NOT IN (
								  SELECT version FROM core
								  ORDER BY version DESC
								  LIMIT 5
							  )`)
		}
	}
}

func main() {
	postgresURL := os.Getenv("POSTGRES")
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		log.Fatal("open", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("ping", err)
	}

	_, err = db.Exec(queryCreateStatesTable)
	if err != nil {
		log.Fatal(err)
	}

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
	timeSave, _ := strconv.Atoi(os.Getenv("TIMELOOP"))
	go saveToDB(db, store, time.Now().UnixMilli(), closeChan, time.Duration(timeSave)*time.Second)

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
