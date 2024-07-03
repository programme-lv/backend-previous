package eval

import (
	"context"
)

type Repository interface {
	AddSubmission(ctx context.Context, submission Submission) error
	ReserveEvaluationID(ctx context.Context) (int64, error)
}
