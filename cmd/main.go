package main

import (
	"2025/internal/server"
	"2025/internal/service"
	"2025/internal/storage"
	"fmt"
	"log"
	"sync"
)

func main() {
	fmt.Println("Test")

	tasks := make(chan service.Task, 100)
	go service.StartWorkerPool(10, tasks)

	jsonPath := "internal\\storage\\storage.json"
	storage, err := storage.NewStorage(jsonPath)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	addr := "localhost:8080"
	server := server.NewServer(addr, storage, tasks)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	// close(tasks)
	wg.Wait()
}
