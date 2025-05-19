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
	ID        string    `json:"id"`         // ID unik untuk identifikasi task
	Data      string    `json:"data"`       // Data atau payload dari task
	Status    string    `json:"status"`     // Status: pending, processing, done, failed
	CreatedAt time.Time `json:"created_at"` // Timestamp saat task dibuat
}

// Global state aplikasi
var (
	queue    = make(chan *Task, 100)  // Channel dengan buffer 100 untuk antrean task
	taskList = make(map[string]*Task) // Map untuk menyimpan semua task berdasarkan ID-nya
	mu       sync.Mutex               // Mutex untuk menjaga konsistensi saat mengakses taskList
)

// generateSimpleID buat ID unik sederhana dari timestamp + angka random
func generateSimpleID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

// writeJSONResponse menulis response JSON dengan status code dan payload
func writeJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// successResponse mengirim response standar sukses dengan message dan optional data
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

// errorResponse mengirim response standar error dengan status code dan pesan error
func errorResponse(w http.ResponseWriter, statusCode int, message string) {
	resp := map[string]interface{}{
		"status":  "error",
		"message": message,
	}
	writeJSONResponse(w, statusCode, resp)
}

// healthHandler adalah endpoint untuk health check service
// Ini digunakan untuk memastikan bahwa service sedang berjalan
func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}
	successResponse(w, "Service is healthy", nil)
}

// enqueueHandler menerima POST request untuk mendaftarkan task baru
// Task akan disimpan dan dimasukkan ke dalam antrean jika valid
func enqueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		errorResponse(w, http.StatusBadRequest, "Invalid task JSON")
		return
	}

	task.Data = strings.TrimSpace(task.Data)
	if task.Data == "" {
		errorResponse(w, http.StatusBadRequest, "data field is required")
		return
	}

	// Generate ID jika kosong
	if strings.TrimSpace(task.ID) == "" {
		task.ID = generateSimpleID()
	} else {
		task.ID = strings.TrimSpace(task.ID)
	}

	// Set status awal dan timestamp
	task.Status = "pending"
	task.CreatedAt = time.Now()

	// Simpan ke map taskList, tapi pastikan tidak ada duplikat ID
	mu.Lock()
	if _, exists := taskList[task.ID]; exists {
		mu.Unlock()
		errorResponse(w, http.StatusBadRequest, "Task ID already exists")
		return
	}
	taskList[task.ID] = &task
	mu.Unlock()

	// Kirim task ke antrean, jika channel penuh dalam 1 detik maka gagal
	select {
	case queue <- &task:
		successResponse(w, "Task queued", map[string]string{"id": task.ID})
	case <-time.After(1 * time.Second):
		// Rollback jika gagal masuk ke antrean
		mu.Lock()
		delete(taskList, task.ID)
		mu.Unlock()
		errorResponse(w, http.StatusServiceUnavailable, "Queue is full, try again later")
	}
}

// listHandler mengembalikan semua task dalam urutan waktu dibuat
// Cocok untuk dashboard admin atau monitoring
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

	// Urutkan berdasarkan CreatedAt agar urut sesuai antrian
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	successResponse(w, "List of tasks", tasks)
}

// retryHandler digunakan untuk me-retry task yang gagal
// Hanya bisa dijalankan pada task dengan status "failed"
func retryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResponse(w, http.StatusMethodNotAllowed, "Only POST allowed")
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 || parts[1] != "retry" {
		errorResponse(w, http.StatusBadRequest, "Bad retry request format. Use /retry/{id}")
		return
	}
	id := parts[2]

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
		// Berhasil retry
	default:
		task.Status = "failed"
		mu.Unlock()
		errorResponse(w, http.StatusServiceUnavailable, "Queue is full, try again later")
		return
	}
	mu.Unlock()

	successResponse(w, "Task retried", map[string]string{"id": task.ID})
}

// worker adalah goroutine yang mengambil task dari queue dan memprosesnya
// Ini adalah simulasi worker backend asynchronous
func worker(ctx context.Context) {
	for {
		select {
		case task := <-queue:
			func() {
				defer func() {
					// Recover jika terjadi panic saat memproses task
					if r := recover(); r != nil {
						log.Printf("Recovered worker panic: %v", r)
						mu.Lock()
						task.Status = "failed"
						mu.Unlock()
					}
				}()

				// Tandai status sebagai sedang diproses
				mu.Lock()
				task.Status = "processing"
				mu.Unlock()

				log.Printf("Processing task: %s", task.ID)
				time.Sleep(2 * time.Second) // Simulasi proses

				// Random failure 20% untuk simulasi error
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

// main adalah entrypoint dari aplikasi
// Men-setup HTTP server, worker, dan menangani graceful shutdown
func main() {
	// Seed untuk random agar hasil failure acak
	rand.Seed(time.Now().UnixNano())

	// Gunakan context untuk mengontrol worker shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Jalankan 1 worker secara paralel
	go worker(ctx)

	// Daftarkan endpoint
	http.HandleFunc("/healthz", healthHandler)
	http.HandleFunc("/enqueue", enqueueHandler)
	http.HandleFunc("/queue", listHandler)
	http.HandleFunc("/retry/", retryHandler) // trailing slash penting

	// Buat HTTP server
	server := &http.Server{Addr: ":8080"}

	// Jalankan goroutine untuk menangani sinyal shutdown
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		sig := <-c
		log.Printf("Got signal %s, shutting down...", sig)
		cancel()

		// Shutdown HTTP server dengan timeout
		ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelTimeout()
		if err := server.Shutdown(ctxTimeout); err != nil {
			log.Fatalf("Server Shutdown Failed:%+v", err)
		}
	}()

	log.Println("Server running on :8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe(): %v", err)
	}

	log.Println("Server exited properly")
}
