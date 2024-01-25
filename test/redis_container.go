package test

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
)

type RedisContainer struct {
	*redis.RedisContainer
	Endpoint string
}

func NewRedisContainer(ctx context.Context) (*RedisContainer, error) {
	container, err := redis.RunContainer(
		ctx,
		testcontainers.WithImage("docker.io/redis:7"),
		redis.WithSnapshotting(10, 1),
		redis.WithLogLevel(redis.LogLevelVerbose),
	)
	if err != nil {
		return nil, err
	}

	endpoint, err := container.Endpoint(ctx, "")
	if err != nil {
		return nil, err
	}

	return &RedisContainer{
		RedisContainer: container,
		Endpoint:       endpoint,
	}, nil
}
