package model

import (
	"time"
)

type Task struct {
	ID        string    `json:"id"`
	Data      string    `json:"data"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
