package worker

import (
	"context"
	"log"

	"github.com/ezzycreative1/simple-queue-system/backend/internal/model"
	"github.com/ezzycreative1/simple-queue-system/backend/internal/service"
)

func StartWorker(ctx context.Context, ts *service.TaskService) {
	go func() {
		for {
			select {
			case task := <-ts.Queue:
				processTask(ts, task)
			case <-ctx.Done():
				log.Println("Worker shutting down...")
				return
			}
		}
	}()
}

func processTask(ts *service.TaskService, task *model.Task) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
			ts.MarkFailed(task)
		}
	}()

	log.Printf("Processing task: %s", task.ID)
	ts.Process(task)

	if task.Status == "failed" {
		log.Printf("Task failed: %s", task.ID)
	} else {
		log.Printf("Task done: %s", task.ID)
	}
}
