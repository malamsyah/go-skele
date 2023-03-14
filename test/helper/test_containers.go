package helper

import (
	"context"
	"fmt"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const RedisClusterLogOccurrence = 3

type imageConfigRequest struct {
	Request    testcontainers.ContainerRequest
	DefaultEnv map[string]string
}

var imageConfigs = map[string]imageConfigRequest{
	"spanner": {
		Request: testcontainers.ContainerRequest{
			Image:        "gcr.io/cloud-spanner-emulator/emulator:latest",
			ExposedPorts: []string{"9010/tcp", "9020/tcp"},
			WaitingFor:   wait.ForLog("gRPC server listening at"),
		},
		DefaultEnv: nil,
	},
	"wiremock": {
		Request: testcontainers.ContainerRequest{
			Image:        "wiremock/wiremock:latest",
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForLog("verbose"),
		},
		DefaultEnv: nil,
	},
	"redis-cluster": {
		Request: testcontainers.ContainerRequest{
			Image:        "grokzen/redis-cluster:6.2.0",
			ExposedPorts: []string{"7000/tcp", "7001/tcp", "7002/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections").WithOccurrence(RedisClusterLogOccurrence),
		},
		DefaultEnv: map[string]string{
			"IP":                "0.0.0.0",
			"INITIAL_PORT":      "7000",
			"MASTERS":           "3",
			"SLAVES_PER_MASTER": "0",
		},
	},
}

type TestContainers struct {
	ctx          context.Context
	containerMap map[string]testcontainers.Container
}

func (tc *TestContainers) Start(name string, image string, envs map[string]string) error {
	config, found := imageConfigs[image]
	if !found {
		return fmt.Errorf("image config for %v is not defined", image)
	}

	containerRequest := config.Request

	_, running := tc.containerMap[name]
	if running {
		return fmt.Errorf("%v container is already running", name)
	}

	defaultEnv := config.DefaultEnv

	for key, val := range envs {
		defaultEnv[key] = val
	}

	containerRequest.Env = defaultEnv
	container, err := testcontainers.GenericContainer(
		tc.ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: containerRequest,
			Started:          true,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to start %v container: %w", name, err)
	}

	tc.containerMap[name] = container

	return nil
}

func (tc *TestContainers) Stop(name string) error {
	container, running := tc.containerMap[name]
	if !running {
		return fmt.Errorf("%v container is not running", name)
	}

	err := container.Terminate(tc.ctx)
	if err != nil {
		return fmt.Errorf("failed to terminate %v container: %w", name, err)
	}

	delete(tc.containerMap, name)

	return nil
}

func (tc *TestContainers) StopAll() error {
	for name, container := range tc.containerMap {
		err := container.Terminate(tc.ctx)
		if err != nil {
			return fmt.Errorf("failed to terminate %v container: %w", name, err)
		}

		delete(tc.containerMap, name)
	}

	return nil
}

func (tc *TestContainers) Host(name string) (string, error) {
	container, running := tc.containerMap[name]
	if !running {
		return "", fmt.Errorf("%v container is not running", name)
	}

	host, err := container.Host(tc.ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get %v container host: %w", name, err)
	}

	return host, nil
}

func (tc *TestContainers) MappedPort(name string, port int) (int, error) {
	container, running := tc.containerMap[name]
	if !running {
		return 0, fmt.Errorf("%v container is not running", name)
	}

	mappedPort, err := container.MappedPort(tc.ctx, nat.Port(fmt.Sprintf("%d/tcp", port)))
	if err != nil {
		return 0, fmt.Errorf("failed to get %v container mapped port %v: %w", name, port, err)
	}

	return mappedPort.Int(), nil
}

func NewTestContainers(ctx context.Context) *TestContainers {
	return &TestContainers{
		ctx:          ctx,
		containerMap: make(map[string]testcontainers.Container),
	}
}
