package main

import (
	"2025/internal/server"
	"2025/internal/service"
	"2025/internal/storage"
	"context"
	"fmt"
	"log"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

func main() {

	// Канал задач для воркеров
	tasks := make(chan service.Task, 100)
	var wg sync.WaitGroup

	go service.StartWorkerPool(10, tasks, &wg)

	jsonPath := filepath.Join("internal", "storage", "storage.json")
	strg, err := storage.NewStorage(jsonPath, nil)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	// Контекст ловит SIGINT / SIGTERM
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// Инициализация сервера
	addr := "localhost:8080"
	srv := server.NewServer(addr, strg, tasks)

	// Запуск сервера
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server failed to start: %v\n", err)
		}
	}()

	// Ожидаем сигнала остановки
	<-ctx.Done()

	// Принять новые задачи ещё 3 секунды
	fmt.Println("Accepting new tasks for 3 more seconds.")
	time.Sleep(3 * time.Second)
	fmt.Println("Starting graceful shutdown.")

	// завершение приёма задач
	close(tasks)

	// завершение воркеров
	wg.Wait()

	// Сохраняем хранилище на диск
	if err := strg.SaveToDisk(); err != nil {
		log.Printf("storage save error: %v", err)
	}

}
