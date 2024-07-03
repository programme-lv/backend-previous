package eval

import "time"

type Evaluation struct {
	id            int64
	evalStatus    string
	totalScore    int64  // received score
	possibleScore *int64 // max score
	compileRData  *runtimeData
	createdAt     time.Time
	taskVersionID int64
}

type runtimeData struct {
	timeMillis      int
	memoryKilobytes int
	exitCode        int
	stdout          string
	stderr          string
}

func NewEvaluation(id int64) *Evaluation {
	return &Evaluation{}
}

func (e *Evaluation) ID() int64 {
	return e.id
}

func (e *Evaluation) Status() string {
	return e.evalStatus
}

func (e *Evaluation) TotalReceivedScore() int64 {
	return e.totalScore
}

func (e *Evaluation) MaxPossibleScore() *int64 {
	return e.possibleScore
}

func (e *Evaluation) CompilationRuntimeData() *runtimeData {
	return e.compileRData
}

func (e *Evaluation) CreatedAt() time.Time {
	return e.createdAt
}

func (e *Evaluation) TaskVersionID() int64 {
	return e.taskVersionID
}

func (rd *runtimeData) TimeMilliseconds() int {
	return rd.timeMillis
}

func (rd *runtimeData) MemoryKilobytes() int {
	return rd.memoryKilobytes
}

func (rd *runtimeData) ExitCode() int {
	return rd.exitCode
}

func (rd *runtimeData) Stdout() string {
	return rd.stdout
}

func (rd *runtimeData) Stderr() string {
	return rd.stderr
}
