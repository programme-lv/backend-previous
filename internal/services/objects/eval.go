package objects

import "time"

type Evaluation struct {
	ID             int64
	StatusID       string
	TotalScore     int64
	PossibleScore  int64
	RuntimeStatsID int64
	CompileDataID  int64
	CreatedAt      time.Time
	TaskVersionID  int64
}

type Submission struct {
	Content    string
	LanguageID string
}

type SubmissionData struct {
	Submission
	ID            int64
	UserID        int64
	TaskID        int64
	CreatedAt     time.Time
	Hidden        bool
	VisibleEvalID int64
}
