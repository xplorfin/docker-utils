package docker

import (
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
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
	client := TestClient{NewDockerClient()}
	defer client.TeardownSession()
	id, err := client.CreateTestContainer()
	if err != nil {
		t.Error(err)
	}
	res, err := client.Exec(id, []string{"echo", "hi"})
	if err != nil {
		t.Error(err)
	}
	if res.ExitCode != 0 && res.StdErr != "" && res.StdOut != "hi" {
		t.Errorf("expeceted output to equal %s", "hi")
	}
}
