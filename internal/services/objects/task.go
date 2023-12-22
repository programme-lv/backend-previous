package objects

import "time"

type TaskVersion struct {
	ID   int64
	Code string
	Name string

	Description *Description

	TimeLimitMs   int
	MemoryLimitKb int

	CreatedAt time.Time
	updatedAt *time.Time
}

type Description struct {
	ID       int64
	Story    string
	Input    string
	Output   string
	Examples []Example
	Notes    string
}

type Example struct {
	ID     int64
	Input  string
	Answer string
}
