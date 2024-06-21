package eval

import (
	"github.com/programme-lv/backend/internal/lang"
	"github.com/programme-lv/backend/internal/task"
	"github.com/programme-lv/backend/internal/user"
	"time"
)

type EvalTestRes struct {
	ID           int64
	EvaluationID int64
	EvalStatusID string
	TaskVTestID  int64
	ExecRData    *RuntimeData
	CheckerRData *RuntimeData
}

type Evaluation struct {
	ID            int64
	TaskVersionID int64

	StatusID      string
	ReceivedScore int64
	PossibleScore *int64

	CheckerRunData *RuntimeData

	TestResults []EvalTestRes

	CreatedAt time.Time
}

type TaskSubmission struct {
	ID          int64
	Author      *user.User
	Language    *lang.ProgrammingLanguage
	Content     string
	Task        *task.Task
	VisibleEval *Evaluation
	Hidden      bool
	CreatedAt   time.Time
}

type RuntimeData struct {
	ID int64

	Stdout *string
	Stderr *string

	TimeMillis      *int64
	MemoryKibibytes *int64
	TimeWallMillis  *int64

	ExitCode *int64
}
