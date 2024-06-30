package eval

import "github.com/google/uuid"

type Submission struct {
	id       int64
	taskID   int64
	authorID int64
	pLangID  string
	msgBody  string
}

func NewSubmission(taskID, authorID int64, pLangID, msgBody string) (*Submission, error) {
	if taskID <= 0 || authorID <= 0 || pLangID == "" {
		return nil, newErrorInvalidSubmissionParams()
	}

	const maxMsgBodyLengthInBytes = 1000 * 256 // 256KB
	if len([]byte(msgBody)) > maxMsgBodyLengthInBytes {
		return nil, newErrorSubmissionBodyTooLarge()
	}

	return &Submission{
		id:       uuid.New(),
		taskID:   taskID,
		authorID: authorID,
		pLangID:  pLangID,
		msgBody:  msgBody,
	}, nil
}

func (s *Submission) ID() int64 {
	return s.id
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
