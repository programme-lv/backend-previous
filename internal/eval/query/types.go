package query

import (
	"github.com/google/uuid"
	"time"
)

type Submission struct {
	UUID             uuid.UUID
	TaskFullName     string
	TaskCode         string
	AuthorUsername   string
	ProgLangID       string
	ProgLangFullName string
	SubmissionCode   string
	EvaluationRes    *Evaluation
	CreatedAt        time.Time
}

type Evaluation struct {
	ID     int64
	Status string

	TotalScore int64
	MaxScore   *int64

	CompileRData *RuntimeData
}

type RuntimeData struct {
	TimeMillis int
	MemoryKB   int
	ExitCode   int
	Stdout     string
	Stderr     string
}
