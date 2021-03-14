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

func (c TestClient) CreateTestContainer() (id string, err error) {
	created, err := c.CreateContainer(&container.Config{
		Image: "docker.io/library/alpine",
		Tty:   false,
	}, nil, nil, nil, "")

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
	output := gofakeit.Word()
	client := TestClient{NewDockerClient()}
	defer func() {
		Nil(t, client.TeardownSession())
	}()
	id, err := client.CreateTestContainer()
	Nil(t, err)

	res, err := client.Exec(id, []string{"echo", output})
	Nil(t, err)

	status, err := client.ContainerStatus(id)
	Nil(t, err)
	Equal(t, status, ContainerRunning)

	if res.ExitCode != 0 && res.StdErr != "" && res.StdOut != output {
		t.Errorf("expeceted output to equal %s", "hi")
	}
}
