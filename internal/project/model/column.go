package model

import (
	"time"
)

// Column represents a Project Column.
type Column struct {
	ID         string    `db:"column_id" json:"id"`
	Title      string    `db:"title" json:"title"`
	ColumnName string    `db:"column_name" json:"columnName"`
	TaskIDS    []string  `db:"task_ids" json:"taskIds"`
	ProjectID  string    `db:"project_id" json:"projectId"`
	UpdatedAt  time.Time `db:"updated_at" json:"updatedAt"`
	CreatedAt  time.Time `db:"created_at" json:"createdAt"`
}

// NewColumn represents a new Column.
type NewColumn struct {
	Title      string `json:"title"`
	ColumnName string `json:"columnName"`
	ProjectID  string `json:"projectId"`
}

// UpdateColumn represents a Column update.
type UpdateColumn struct {
	Title     *string   `json:"title"`
	TaskIDS   *[]string `json:"taskIds"`
	UpdatedAt time.Time `json:"updatedAt"`
}
