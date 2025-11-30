package service

import (
	"2025/internal/check"
	"log"
)

type Result struct {
	URL    string
	Status string
}

type Task struct {
	URL string
	Res chan Result
}

// StartWorkerPool запускает N воркеров,
// которые будут обрабатывать задачи из канала tasks
func StartWorkerPool(n int, tasks chan Task) {
	for i := 1; i <= n; i++ {
		go func(workerId int) {
			// Воркеры читают задачи из канала, пока канал не будет закрыт
			for task := range tasks {
				log.Printf("worker %d processing %s", workerId, task.URL)

				if check.CheckLink(task.URL) {
					task.Res <- Result{URL: task.URL, Status: "available"}
				} else {
					task.Res <- Result{URL: task.URL, Status: "not available"}
				}
			}
		}(i)
	}
}
