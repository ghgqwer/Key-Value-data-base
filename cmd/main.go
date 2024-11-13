package main

import (
	"BolshoiGolangProject/internal/filework"
	"BolshoiGolangProject/internal/server"
	"BolshoiGolangProject/internal/storage/storage"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "BolshoiGolangProject/docs"

	_ "github.com/lib/pq"
)

const (
	queryCreateStatesTable = `CREATE TABLE IF NOT EXISTS core (
		version bigserial PRIMARY KEY,
		timestamp bigint NOT NULL,
		payload JSONB NOT NULL
	)`

	queryCreateState = `INSERT INTO core (timestamp, payload) VALUES ($1, $2)`

	execVacuum = `VACUUM core`

	execVacuumFull = `VACUUM FULL core`
)

const (
	postgresEnv   = "POSTGRES"
	serverPortEnv = "BASIC_SERVER_PORT"
	timeLoopEnv   = "TIMELOOP"
)

type Core struct {
	Version   int    `json: version`
	Timestamp int    `json: timestamp`
	Payload   string `json: payload`
}

func vacuumDB(db *sql.DB, closeChan <-chan struct{}, n time.Duration) {
	for {
		select {
		case <-closeChan:
			return
		case <-time.After(n):
			db.Exec(execVacuum)
		}
	}

}

func fullVacuumDB(db *sql.DB, closeChan <-chan struct{}, n time.Duration) {
	for {
		select {
		case <-closeChan:
			return
		case <-time.After(n):
			db.Exec(execVacuumFull)
		}
	}

}

func saveToDB(db *sql.DB, storage storage.Storage, timestamp int64, closeChan <-chan struct{}, n time.Duration) {
	for {
		select {
		case <-closeChan:
			return
		case <-time.After(n):
			payload, _ := json.Marshal(storage)
			db.Exec(queryCreateState, timestamp, payload)

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
	if err := runMain(); err != nil {
		log.Fatal()
	}
}

func runMain() error {
	postgresURL := os.Getenv(postgresEnv)
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	_, err = db.Exec(queryCreateStatesTable)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	//serverPort := "8090"

	serverPort, ok := os.LookupEnv(serverPortEnv)
	if !ok {
		fmt.Println("not port provided")
		os.Exit(1)
	}

	store, err := storage.NewStorage()
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %v", err)
	}

	if err := filework.ReadFromJSON(store, server.DataJson); err != nil {
		return fmt.Errorf("failed to read from JSON: %v", err)
	}

	closeChan := make(chan struct{})
	go store.GarbageCollection(closeChan, 10*time.Second)
	go store.LoggerSync(closeChan, 10*time.Second)
	go vacuumDB(db, closeChan, 100*time.Second)
	go fullVacuumDB(db, closeChan, 1000*time.Second)
	timeSave, _ := strconv.Atoi(os.Getenv(timeLoopEnv))
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

	return nil
}
