package server

import (
	"net/http"

	"github.com/ezzycreative1/simple-queue-system/backend/internal/handlers"
	"github.com/ezzycreative1/simple-queue-system/backend/internal/service"
)

func RegisterRoutes(mux *http.ServeMux, ts *service.TaskService) {
	mux.HandleFunc("/api/healthz", handlers.HealthHandler())
	mux.HandleFunc("/api/enqueue", handlers.EnqueueHandler(ts))
	mux.HandleFunc("/api/queue", handlers.ListHandler(ts))
	mux.HandleFunc("/api/retry/", handlers.RetryHandler(ts))
}
