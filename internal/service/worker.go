package service

import "2025/internal/check"

type Task struct {
	URL    string
	Result chan string
}

// StartWorkerPool запускает N воркеров,
// которые будут обрабатывать задачи из канала tasks
func StartWorkerPool(n int, tasks chan Task) {
	for i := 1; i <= n; i++ {
		go func(workerId int) {
			// Воркеры читают задачи из канала, пока канал не будет закрыт
			for task := range tasks {
				if check.CheckLink(task.URL) {
					task.Result <- "available"
				} else {
					task.Result <- "not available"
				}
			}
		}(i)
	}
}
