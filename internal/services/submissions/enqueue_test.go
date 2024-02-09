package submissions

import (
	"os"
	"testing"

	"github.com/programme-lv/backend/internal/services/objects"
	"github.com/programme-lv/backend/internal/testing/testrmq"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

var (
	rmq *amqp.Connection
	// db  *sqlx.DB
)

func TestMain(m *testing.M) {
	// dbContainer, err := testdb.NewMigratedPostgresTestcontainer()
	// if err != nil {
	// 	panic(err)
	// }
	// db = dbContainer.GetConn()
	// defer dbContainer.Close()

	rmqContainer, err := testrmq.NewRMQTestcontainer()
	if err != nil {
		panic(err)
	}
	rmq = rmqContainer.GetConn()
	defer rmqContainer.Close()

	code := m.Run()
	os.Exit(code)
}

func TestEnqueueEvaluationIntoRMQ(t *testing.T) {
	submission := objects.RawSubmission{
		Content:    "print('Hello World')",
		LanguageID: "python",
	}

	eval := objects.EvaluationJob{
		ID:            69,
		TaskVersionID: 420,
	}
	err := EnqueueEvaluationIntoRMQ(rmq, submission, eval)
	if err != nil {
		t.Errorf("EnqueueEvaluationIntoRMQ() error = %v", err)
	}

	recSubm, recEval, err := DequeueEvaluationFromRMQ(rmq)
	if err != nil {
		t.Errorf("DequeueEvaluationFromRMQ() error = %v", err)
	}

	assert.Equal(t, submission, *recSubm)
	assert.Equal(t, eval, *recEval)

	t.Logf("recSubm: %+v", recSubm)
	t.Logf("recEval: %+v", recEval)

}
