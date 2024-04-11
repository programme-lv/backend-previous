package objects

import "time"

type TaskVersion struct {
	ID     int64
	TaskID int64
	Code   string
	Name   string

	Description *Description

	TimeLimitMs   int64
	MemoryLimitKb int64

	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Task struct {
	ID          int64
	CreatedByID int64

	Current *TaskVersion
	Stable  *TaskVersion

	CreatedAt time.Time
}

type Description struct {
	ID       int64
	Story    string
	Input    string
	Output   string
	Examples []Example
	Notes    *string
}

type Example struct {
	ID     int64
	Input  string
	Answer string
}
