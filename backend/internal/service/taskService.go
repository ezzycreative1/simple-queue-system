package service

import (
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/ezzycreative1/simple-queue-system/backend/internal/model"
)

type TaskService struct {
	Queue    chan *model.Task
	taskList map[string]*model.Task
	mu       sync.Mutex
}

func NewTaskService(bufferSize int) *TaskService {
	return &TaskService{
		Queue:    make(chan *model.Task, bufferSize),
		taskList: make(map[string]*model.Task),
	}
}

func (s *TaskService) generateSimpleID() string {
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

func (s *TaskService) AddTask(id, data string) (*model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id == "" {
		id = s.generateSimpleID()
	} else {
		if _, exists := s.taskList[id]; exists {
			return nil, errors.New("Task ID already exists")
		}
	}

	task := &model.Task{
		ID:        id,
		Data:      data,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.taskList[task.ID] = task

	select {
	case s.Queue <- task:
		return task, nil
	case <-time.After(1 * time.Second):
		delete(s.taskList, task.ID)
		return nil, errors.New("Queue is full, try again later")
	}
}

func (s *TaskService) ListTasks(status string, page, limit int) ([]*model.Task, int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := make([]*model.Task, 0, len(s.taskList))
	for _, t := range s.taskList {
		if status == "" || t.Status == status {
			filtered = append(filtered, t)
		}
	}

	// Sort by CreatedAt ascending
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
	})

	total := len(filtered)
	start := (page - 1) * limit
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	return filtered[start:end], total
}

func (s *TaskService) RetryTask(id string) (*model.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.taskList[id]
	if !exists {
		return nil, errors.New("Task not found")
	}
	if task.Status != "failed" {
		return nil, errors.New("Only failed tasks can be retried")
	}

	task.Status = "pending"
	task.UpdatedAt = time.Now()

	select {
	case s.Queue <- task:
		return task, nil
	default:
		task.Status = "failed"
		return nil, errors.New("Queue is full, try again later")
	}
}

func (s *TaskService) Process(task *model.Task) {
	s.mu.Lock()
	task.Status = "processing"
	task.UpdatedAt = time.Now()
	s.mu.Unlock()

	// Simulate processing
	time.Sleep(2 * time.Second)

	s.mu.Lock()
	defer s.mu.Unlock()
	if rand.Intn(5) == 0 {
		task.Status = "failed"
	} else {
		task.Status = "done"
	}
	task.UpdatedAt = time.Now()
}

func (s *TaskService) MarkFailed(task *model.Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	task.Status = "failed"
	task.UpdatedAt = time.Now()
}
