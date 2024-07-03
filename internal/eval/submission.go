package eval

import "github.com/google/uuid"

type Submission struct {
	uuid      uuid.UUID
	taskID    int64
	authorID  int64
	pLangID   string
	msgBody   string
	evals     []*Evaluation
	relevEval *Evaluation // relevant / visible evaluation
}

func NewSubmission(uuid uuid.UUID, taskID, authorID int64, pLangID, msgBody string) (*Submission, error) {
	if taskID <= 0 || authorID <= 0 || pLangID == "" {
		return nil, newErrorInvalidSubmissionParams()
	}

	const maxMsgBodyLengthInBytes = 1000 * 256 // 256KB
	if len([]byte(msgBody)) > maxMsgBodyLengthInBytes {
		return nil, newErrorSubmissionBodyTooLarge()
	}

	return &Submission{
		uuid:     uuid,
		taskID:   taskID,
		authorID: authorID,
		pLangID:  pLangID,
		msgBody:  msgBody,
		evals:    make([]*Evaluation, 0),
	}, nil
}

func (s *Submission) Evaluate(evaluationID int64) {
	evaluation := NewEvaluation(evaluationID)
}

func (s *Submission) UUID() uuid.UUID {
	return s.uuid
}

func (s *Submission) TaskID() int64 {
	return s.taskID
}

func (s *Submission) AuthorID() int64 {
	return s.authorID
}

func (s *Submission) ProgrammingLanguageID() string {
	return s.pLangID
}

func (s *Submission) MessageBody() string {
	return s.msgBody
}

func (s *Submission) EvaluationResult() *Evaluation {
	return s.relevEval
}
