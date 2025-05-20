package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ezzycreative1/simple-queue-system/backend/internal/model"
	"github.com/ezzycreative1/simple-queue-system/backend/internal/service"
	"github.com/ezzycreative1/simple-queue-system/backend/pkg/util"
)

func HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
			return
		}
		util.SuccessResponse(w, "Service is healthy", nil)
	}
}

func EnqueueHandler(ts *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST allowed")
			return
		}

		var payload model.EnqueueRequest
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			util.ErrorResponse(w, http.StatusBadRequest, "Invalid task JSON")
			return
		}

		payload.Data = strings.TrimSpace(payload.Data)
		if payload.Data == "" {
			util.ErrorResponse(w, http.StatusBadRequest, "data field is required")
			return
		}

		task, err := ts.AddTask(strings.TrimSpace(payload.ID), payload.Data)
		if err != nil {
			util.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		util.SuccessResponse(w, "Task queued", map[string]string{"id": task.ID})
	}
}

func ListHandler(ts *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only GET allowed")
			return
		}

		status := strings.TrimSpace(r.URL.Query().Get("status"))
		pageStr := r.URL.Query().Get("page")
		limitStr := r.URL.Query().Get("limit")

		page := 1
		limit := 20
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}

		tasks, total := ts.ListTasks(status, page, limit)

		util.SuccessResponse(w, "List of tasks", map[string]interface{}{
			"tasks": tasks,
			"meta": map[string]interface{}{
				"page":  page,
				"limit": limit,
				"total": total,
				"count": len(tasks),
			},
		})
	}
}

func RetryHandler(ts *service.TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			util.ErrorResponse(w, http.StatusMethodNotAllowed, "Only POST allowed")
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/api/retry/")
		id = strings.TrimSpace(id)
		if id == "" {
			util.ErrorResponse(w, http.StatusBadRequest, "Task ID is required in URL")
			return
		}

		task, err := ts.RetryTask(id)
		if err != nil {
			util.ErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		util.SuccessResponse(w, "Task retried", map[string]string{"id": task.ID})
	}
}
