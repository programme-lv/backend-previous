package objects

import "time"

type EvaluationJob struct {
	ID            int64
	TaskVersionID int64
}

type Evaluation struct {
	ID            int64
	TaskVersionID int64

	StatusID       string
	TotalScore     int64
	PossibleScore  int64
	RuntimeStatsID int64
	CompileDataID  int64
	CreatedAt      time.Time
}

type RawSubmission struct {
	Content    string
	LanguageID string
}

type TaskSubmission struct {
	RawSubmission
	ID            int64
	UserID        int64
	TaskID        int64
	CreatedAt     time.Time
	Hidden        bool
	VisibleEvalID int64
}
