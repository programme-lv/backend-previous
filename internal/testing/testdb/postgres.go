package testdb

import (
	"context"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type postgresContainer struct {
	container testcontainers.Container
	user      string
	password  string
	database  string
}

func startPostgresContainer(networkName string, networkAlias string, user string, pass string, dbName string) (*postgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     user,
			"POSTGRES_PASSWORD": pass,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForExposedPort(),
		),
		Networks: []string{networkName},
		NetworkAliases: map[string][]string{
			networkName: {networkAlias},
		},
		Name: networkAlias,
	}

	container, err := testcontainers.GenericContainer(context.Background(),
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
	if err != nil {
		return nil, err
	}

	res := &postgresContainer{
		container: container,
		user:      user,
		password:  pass,
		database:  dbName,
	}
	return res, nil
}
