package docker

import (
	"testing"

	"github.com/brianvoe/gofakeit/v5"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	. "github.com/stretchr/testify/assert"
)

type TestClient struct {
	Client
}

func (c TestClient) CreateTestContainer(network string) (id string, err error) {
	created, err := c.CreateContainer(&container.Config{
		Image: "docker.io/library/alpine",
		Tty:   false,
	}, &container.HostConfig{
		NetworkMode: container.NetworkMode(network),
	}, nil, nil, "")

	if err != nil {
		return id, err
	}

	err = c.ContainerStart(c.Ctx, created.ID, types.ContainerStartOptions{})
	if err != nil {
		return id, err
	}
	c.PrintContainerLogs(created.ID)
	return created.ID, err
}

func TestContainerLifecycle(t *testing.T) {
	client := TestClient{NewDockerClient()}
	// we run this twice to ensure teardown occurred correctly
	containerNetworkName := gofakeit.Word()
	containerOutput := gofakeit.Word()

	networkID, err := client.CreateNetwork(containerNetworkName)
	Nil(t, err)

	id, err := client.CreateTestContainer(networkID)
	Nil(t, err)

	res, err := client.Exec(id, []string{"echo", containerOutput})
	Nil(t, err)

	if res.ExitCode != 0 && res.StdErr != "" && res.StdOut != containerOutput {
		t.Errorf("expeceted containerOutput to equal %s", "hi")
	}

	// TODO verify
	_, err = client.ContainerStatus(id)
	Nil(t, err)

	Nil(t, client.TeardownSession())
}
