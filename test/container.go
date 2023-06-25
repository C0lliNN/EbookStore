package test

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"path/filepath"
)

type PostgresContainer struct {
	testcontainers.Container
	URI string
}

func NewPostgresContainer(ctx context.Context) (*PostgresContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.2",
		ExposedPorts: []string{"5432/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections"),
		),
		Env: map[string]string{
			"POSTGRES_DB":       "postgres",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
		},
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	addr, err := container.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("postgres://postgres:postgres@%s/postgres?sslmode=disable", addr)

	return &PostgresContainer{
		Container: container,
		URI:       uri,
	}, nil
}

type LocalstackContainer struct {
	testcontainers.Container
	Port string
}

func NewLocalstackContainer(ctx context.Context) (*LocalstackContainer, error) {
	setupPath, err := filepath.Abs("../../../scripts/setup_localstack.sh")
	if err != nil {
		return nil, err
	}

	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack:0.14.2",
		ExposedPorts: []string{"4566/tcp"},
		WaitingFor:   wait.ForLog("Initialization has finished!"),
		Env: map[string]string{
			"DEFAULT_REGION": "us-east-2",
			"SERVICES":       "ses,s3",
			"START_WEB":      "0",
		},
		Mounts: testcontainers.Mounts(testcontainers.BindMount(setupPath, "/docker-entrypoint-initaws.d/init.sh")),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "4566")
	if err != nil {
		return nil, err
	}

	return &LocalstackContainer{
		Container: container,
		Port:      mappedPort.Port(),
	}, nil
}
