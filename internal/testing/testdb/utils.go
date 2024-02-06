package testdb

import (
	"context"
	"math/rand"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

func extractTestcontainerExternalHostAndPort(container testcontainers.Container) (host string, port string, err error) {
	host, err = container.Host(context.Background())
	if err != nil {
		return
	}

	natPort, err := container.MappedPort(context.Background(), "5432")
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

func createNewNetwork() (*testcontainers.DockerNetwork, error) {
	return network.New(context.Background(), network.WithCheckDuplicate())
}

func randomLowercaseLetterString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
