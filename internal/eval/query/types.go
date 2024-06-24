package query

import "time"

type Submission struct {
	ID               int64
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
