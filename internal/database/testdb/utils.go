package testdb

import (
	"context"
	"math/rand"

	"github.com/testcontainers/testcontainers-go"
)

func extractTestcontainerHostAndPort(container testcontainers.Container) (host string, port string, err error) {
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

func createNetwork(networkName string) (testcontainers.Network, error) {
	network, err := testcontainers.GenericNetwork(context.Background(),
		testcontainers.GenericNetworkRequest{
			NetworkRequest: testcontainers.NetworkRequest{
				Name:           networkName,
				CheckDuplicate: true,
			},
		})
	return network, err
}

func randomLowercaseLetterString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
