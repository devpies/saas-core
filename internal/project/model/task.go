package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var taskValidator *validator.Validate

func init() {
	v := NewValidator()
	taskValidator = v
}

// Task represents a Project Task.
type Task struct {
	ID          string    `db:"task_id" json:"id"`
	Key         string    `db:"key" json:"key"`
	Title       string    `db:"title" json:"title"`
	TenantID    string    `db:"tenant_id" json:"tenantId"`
	Points      int       `db:"points" json:"points"`
	UserID      string    `db:"user_id" json:"userId"`
	Content     string    `db:"content" json:"content"`
	ProjectID   string    `db:"project_id" json:"projectId"`
	AssignedTo  string    `db:"assigned_to" json:"assignedTo"`
	Attachments []string  `db:"attachments" json:"attachments"`
	Comments    []string  `db:"comments" json:"comments"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
}

// NewTask represents a new Task.
type NewTask struct {
	Title string `json:"title" validate:"required,min=1,max=75"`
}

// Validate validates a NewTask.
func (nt *NewTask) Validate() error {
	return taskValidator.Struct(nt)
}

// UpdateTask represents a Task being updated.
type UpdateTask struct {
	Title       *string   `json:"title" validate:"omitempty,max=75"`
	Points      *int      `json:"points"`
	Content     *string   `json:"content" validate:"omitempty,max=1000"`
	AssignedTo  *string   `json:"assignedTo"`
	Attachments []string  `json:"attachments"`
	Comments    []string  `json:"comments"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Validate validates an UpdateTask payload.
func (ut *UpdateTask) Validate() error {
	return taskValidator.Struct(ut)
}

// MoveTask represents a Task being moved between columns.
type MoveTask struct {
	To      string   `json:"to" validate:"required"`
	From    string   `json:"from" validate:"required"`
	TaskIds []string `json:"taskIds" validate:"required"`
}

// Validate validates a MoveTask payload.
func (mt *MoveTask) Validate() error {
	return taskValidator.Struct(mt)
}
