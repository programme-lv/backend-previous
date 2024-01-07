package submissions

import (
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/programme-lv/backend/internal/testing/testdb"
	"github.com/programme-lv/backend/internal/testing/testrmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	rmq *amqp.Connection
	db  *sqlx.DB
)

func TestMain(m *testing.M) {
	dbContainer, err := testdb.NewMigratedPostgresTestcontainer()
	if err != nil {
		panic(err)
	}
	db = dbContainer.GetConn()
	defer dbContainer.Close()

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
	// TODO: implement
	t.Error("not implemented")
}
