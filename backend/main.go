package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Struct Task mewakili unit pekerjaan yang akan dikerjakan oleh worker
type Task struct {
	ID        string    `json:"id"`
	Data      string    `json:"data"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	queue    = make(chan *Task, 100)
	taskList = make(map[string]*Task)
	mu       sync.Mutex
)

func generateSimpleID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func successResponse(w http.ResponseWriter, message string, data interface{}) {
	resp := map[string]interface{}{
		"status":  "success",
		"message": message,
	}
	if data != nil {
		resp["data"] = data
	}
	writeJSONResponse(w, http.StatusOK, resp)
}

func errorResponse(w http.ResponseWriter, statusCode int, message string) {
	resp := map[string]interface{}{
		"status":  "error",
		"message": message,
	}
	writeJSONResponse(w, statusCode, resp)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	successResponse(w, "Service is healthy", nil)
}

func enqueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	type taskPayload struct {
		ID   string `json:"id"`
		Data string `json:"data"`
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	var task taskPayload
	errCh := make(chan error, 1)
	go func() {
		errCh <- json.NewDecoder(r.Body).Decode(&task)
	}()

	select {
	case err := <-errCh:
		if err != nil {
			errorResponse(w, http.StatusBadRequest, "Invalid task JSON")
			return
		}
	case <-ctx.Done():
		errorResponse(w, http.StatusRequestTimeout, "Request timeout")
		return
	}

	task.Data = strings.TrimSpace(task.Data)
	if task.Data == "" {
		errorResponse(w, http.StatusBadRequest, "data field is required")
		return
	}

	id := strings.TrimSpace(task.ID)
	if id == "" {
		id = generateSimpleID()
	}

	newTask := &Task{
		ID:        id,
		Data:      task.Data,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	mu.Lock()
	if _, exists := taskList[newTask.ID]; exists {
		mu.Unlock()
		errorResponse(w, http.StatusBadRequest, "Task ID already exists")
		return
	}
	taskList[newTask.ID] = newTask
	mu.Unlock()

	select {
	case queue <- newTask:
		successResponse(w, "Task queued", map[string]string{"id": newTask.ID})
	case <-time.After(1 * time.Second):
		mu.Lock()
		delete(taskList, newTask.ID)
		mu.Unlock()
		errorResponse(w, http.StatusServiceUnavailable, "Queue is full, try again later")
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResponse(w, http.StatusMethodNotAllowed, "Only GET allowed")
		return
	}

	mu.Lock()
	tasks := make([]*Task, 0, len(taskList))
	for _, t := range taskList {
		tasks = append(tasks, t)
	}
	mu.Unlock()

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	successResponse(w, "List of tasks", tasks)
}

func retryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/retry/")
	id = strings.TrimSpace(id)
	if id == "" {
		errorResponse(w, http.StatusBadRequest, "Task ID is required in URL")
		return
	}

	mu.Lock()
	task, exists := taskList[id]
	if !exists {
		mu.Unlock()
		errorResponse(w, http.StatusNotFound, "Task not found")
		return
	}

	if task.Status != "failed" {
		mu.Unlock()
		errorResponse(w, http.StatusBadRequest, "Only failed tasks can be retried")
		return
	}

	task.Status = "pending"
	select {
	case queue <- task:
		// retried
	default:
		task.Status = "failed"
		mu.Unlock()
		errorResponse(w, http.StatusServiceUnavailable, "Queue is full, try again later")
		return
	}
	mu.Unlock()

	successResponse(w, "Task retried", map[string]string{"id": task.ID})
}

func worker(ctx context.Context) {
	for {
		select {
		case task := <-queue:
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered worker panic: %v", r)
						mu.Lock()
						task.Status = "failed"
						mu.Unlock()
					}
				}()

				mu.Lock()
				task.Status = "processing"
				mu.Unlock()

				log.Printf("Processing task: %s", task.ID)
				time.Sleep(2 * time.Second)

				mu.Lock()
				if rand.Intn(5) == 0 {
					task.Status = "failed"
					log.Printf("Task failed: %s", task.ID)
				} else {
					task.Status = "done"
					log.Printf("Task done: %s", task.ID)
				}
				mu.Unlock()
			}()
		case <-ctx.Done():
			log.Println("Worker shutting down...")
			return
		}
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker(ctx)

	http.HandleFunc("/api/healthz", healthHandler)
	http.HandleFunc("/api/enqueue", enqueueHandler)
	http.HandleFunc("/api/queue", listHandler)
	http.HandleFunc("/api/retry/", retryHandler)

	server := &http.Server{
		Addr:    ":8081",
		Handler: corsMiddleware(http.DefaultServeMux),
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		log.Printf("Got signal %s, shutting down...", sig)
		cancel()

		ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelTimeout()
		if err := server.Shutdown(ctxTimeout); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}
	}()

	log.Println("Server running on :8081")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}
	log.Println("Server exited properly")
}
