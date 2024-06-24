package graphql

import (
	"fmt"
	"github.com/programme-lv/backend/internal/eval/query"
)

func mapSubmissionQueryToGQL(submission query.Submission) (*Submission, error) {
	convInt64ToIntPointer := func(i *int64) *int {
		if i == nil {
			return nil
		}
		convInt := int(*i)
		return &convInt
	}

	var evalResults *Evaluation = nil
	if submission.EvaluationRes != nil {
		evalResults = &Evaluation{
			ID:            fmt.Sprint(submission.EvaluationRes.ID),
			Status:        submission.EvaluationRes.Status,
			TotalScore:    int(submission.EvaluationRes.TotalScore),
			PossibleScore: convInt64ToIntPointer(submission.EvaluationRes.MaxScore),
			CompileRData:  nil,
			TestResults:   nil, // TODO
		}
		if submission.EvaluationRes.CompileRData != nil {
			evalResults.CompileRData = &RuntimeData{
				TimeMs:   submission.EvaluationRes.CompileRData.TimeMillis,
				MemoryKb: submission.EvaluationRes.CompileRData.MemoryKB,
				ExitCode: submission.EvaluationRes.CompileRData.ExitCode,
				Stdout:   submission.EvaluationRes.CompileRData.Stdout,
				Stderr:   submission.EvaluationRes.CompileRData.Stderr,
			}
		}

	}
	marshalledTime, err := submission.CreatedAt.MarshalText()
	if err != nil {
		return nil, err
	}
	return &Submission{
		ID:               fmt.Sprint(submission.ID),
		TaskFullName:     submission.TaskFullName,
		TaskCode:         submission.TaskCode,
		AuthorUsername:   submission.AuthorUsername,
		ProgLangID:       fmt.Sprint(submission.ProgLangID),
		ProgLangFullName: fmt.Sprint(submission.ProgLangFullName),
		SubmissionCode:   submission.SubmissionCode,
		EvalResults:      evalResults,
		CreatedAt:        string(marshalledTime),
	}, nil
}
