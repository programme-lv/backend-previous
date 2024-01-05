package submissions

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/database/testdb"
	"github.com/programme-lv/backend/internal/services/objects"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rmq *amqp.Connection
	db  *sqlx.DB
)

func TestMain(m *testing.M) {
	provider, err := testdb.NewPostgresTestcontainer()
	if err != nil {
		panic(err)
	}
	db = provider.GetConn()
	defer provider.Close()

	code := m.Run()
	os.Exit(code)
}

func TestEnqueueEvaluationIntoRMQ(t *testing.T) {
	type args struct {
		rmq        *amqp.Connection
		submission objects.RawSubmission
		eval       objects.EvaluationJob
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EnqueueEvaluationIntoRMQ(tt.args.rmq, tt.args.submission, tt.args.eval); (err != nil) != tt.wantErr {
				t.Errorf("EnqueueEvaluationIntoRMQ() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
