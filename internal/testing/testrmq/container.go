package testrmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type RMQTestcontainer interface {
	GetConn() *amqp.Connection
	Close()
}

type rmqTestcontainer struct {
	container testcontainers.Container
	conn      *amqp.Connection
}

func (r *rmqTestcontainer) GetConn() *amqp.Connection {
	return nil
}

func (r *rmqTestcontainer) Close() {
	r.conn.Close()
	r.container.Terminate(context.Background())
}

func NewRMQTestcontainer() (RMQTestcontainer, error) {
	username := "guest"
	password := "guest"

	container, err := startRMQContainer(username, password)
	if err != nil {
		return nil, err
	}

	host, port, err := extractTestcontainerExternalHostAndPort(container)
	if err != nil {
		return nil, err
	}

	amqpConnStr := fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port)
	conn, err := amqp.Dial(amqpConnStr)
	if err != nil {
		return nil, err
	}

	return &rmqTestcontainer{
		container: container,
		conn:      conn,
	}, nil

}

func startRMQContainer(defaultUser, defaultPass string) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image: "rabbitmq:3.12",
		ExposedPorts: []string{
			"5672/tcp",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("Server startup complete"),
			wait.ForExposedPort(),
		),
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER": "guest",
			"RABBITMQ_DEFAULT_PASS": "guest",
		},
	}

	container, err := testcontainers.GenericContainer(context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	return container, nil
}

func extractTestcontainerExternalHostAndPort(container testcontainers.Container) (host string, port string, err error) {
	host, err = container.Host(context.Background())
	if err != nil {
		return
	}

	natPort, err := container.MappedPort(context.Background(), "5672")
	if err != nil {
		return
	}

	port = natPort.Port()

	// strip /tcp suffix from port
	if port[len(port)-4:] == "/tcp" {
		port = port[:len(port)-4]
	}

	return
}
