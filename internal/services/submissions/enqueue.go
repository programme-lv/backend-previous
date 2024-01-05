package submissions

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/go-jet/jet/v2/qrm"
	"github.com/programme-lv/backend/internal/services/objects"
	"github.com/programme-lv/tester/pkg/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

func EnqueueEvaluationForTaskSubmission(
	db qrm.Queryable, evaluationID, taskVersionID int64, submission objects.Submission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := messaging.EvaluationRequest{
		TaskVersionId: int64(taskVersionID),
		Submission: messaging.Submission{
			SourceCode: submission.Content,
			LanguageId: submission.LanguageID,
		},
	}

	correlation := messaging.Correlation{
		HasEvaluationId: true,
		EvaluationId:    evaluation.ID,
		UnixMillis:      time.Now().UnixMilli(),
		RandomInt63:     rand.Int63(),
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	correlationJson, err := json.Marshal(correlation)
	if err != nil {
		return err
	}

	ch, err := r.SubmissionRMQ.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("eval_q", true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:   "application/json",
		Body:          bodyJson,
		ReplyTo:       "res_q",
		CorrelationId: string(correlationJson),
	})
	if err != nil {
		return err
	}

	return nil
}
