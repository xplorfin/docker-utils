package docker

import "testing"

func TestPullImage(t *testing.T) {
	client := NewDockerClient()
	err := client.PullImage("docker.io/library/alpine")
	if err != nil {
		t.Error(err)
	}
}
