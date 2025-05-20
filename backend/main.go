package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ezzycreative1/simple-queue-system/backend/internal/middleware"
	"github.com/ezzycreative1/simple-queue-system/backend/internal/service"
	"github.com/ezzycreative1/simple-queue-system/backend/internal/worker"
	"github.com/ezzycreative1/simple-queue-system/backend/server"
)

func main() {
	taskService := service.NewTaskService(100)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	worker.StartWorker(ctx, taskService)

	mux := http.NewServeMux()
	server.RegisterRoutes(mux, taskService)

	server := &http.Server{
		Addr:    ":8081",
		Handler: middleware.CORSMiddleware(mux),
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
