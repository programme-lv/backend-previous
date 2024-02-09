package submissions

import (
	"encoding/json"

	"github.com/programme-lv/backend/internal/services/objects"
	"github.com/programme-lv/tester/pkg/messaging"
	"github.com/rabbitmq/amqp091-go"
)

func DequeueEvaluationFromRMQ(rmq *amqp091.Connection) (
	*objects.RawSubmission, *objects.EvaluationJob, error) {

	ch, err := rmq.Channel()
	if err != nil {
		return nil, nil, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(EvalQueueName, true, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	msg := <-msgs

	var body messaging.EvaluationRequest
	err = json.Unmarshal(msg.Body, &body)
	if err != nil {
		return nil, nil, err
	}

	var correlation messaging.Correlation
	err = json.Unmarshal([]byte(msg.CorrelationId), &correlation)
	if err != nil {
		return nil, nil, err
	}

	submission := objects.RawSubmission{
		Content:    body.Submission.SourceCode,
		LanguageID: body.Submission.LanguageId,
	}

	eval := objects.EvaluationJob{
		ID:            correlation.EvaluationId,
		TaskVersionID: body.TaskVersionId,
	}

	return &submission, &eval, nil
}
