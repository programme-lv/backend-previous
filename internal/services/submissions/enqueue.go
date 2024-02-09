package submissions

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	"github.com/programme-lv/backend/internal/services/objects"
	"github.com/programme-lv/tester/pkg/messaging"
	amqp "github.com/rabbitmq/amqp091-go"
)

func EnqueueEvaluationIntoRMQ(rmq *amqp.Connection,
	submission objects.RawSubmission, eval objects.EvaluationJob) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := messaging.EvaluationRequest{
		TaskVersionId: eval.TaskVersionID,
		Submission: messaging.Submission{
			SourceCode: submission.Content,
			LanguageId: submission.LanguageID,
		},
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return err
	}

	correlation := messaging.Correlation{
		HasEvaluationId: true,
		EvaluationId:    eval.ID,
		UnixMillis:      time.Now().UnixMilli(),
		RandomInt63:     rand.Int63(),
	}

	correlationJson, err := json.Marshal(correlation)
	if err != nil {
		return err
	}

	ch, err := rmq.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(EvalQueueName, true, false, false, false,
		amqp.Table{})
	if err != nil {
		return err
	}

	err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:   "application/json",
		Body:          bodyJson,
		ReplyTo:       ResponseQueueName,
		CorrelationId: string(correlationJson),
	})

	if err != nil {
		return err
	}

	return nil
}
