package domain

import "time"

// TaskVersion is a snapshot of task development state at a certain point in time.
// At the persistence layer, it is supposed to be read-only, i.e. it should not be modified after creation.
// Modifications should be done by creating a new version. A version is considered
// to be successor of another version if it has the same TaskID and its ID is greater.
type TaskVersion struct {
	ID     int64
	TaskID int64
	Code   string
	Name   string

	Statement *Statement

	TimeLimitMs   int64
	MemoryLimitKb int64

	// CreatedAt  is the time when the new version was created, i.e. the time when the task was updated.
	CreatedAt time.Time
}

// Task represents a task that can be solved by a user. It is a collection of task versions.
type Task struct {
	ID      int64
	OwnerID int64

	// Current is the newest / latest version of the task.
	// Accessible only by the creator / owner of the task.
	Current *TaskVersion

	Stable *TaskVersion

	CreatedAt time.Time
}

// Statement is a set of sections that describe the task. The sections are formatted in Markdown.
// Notes are an optional field that can be used to provide additional information to the task.
type Statement struct {
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
