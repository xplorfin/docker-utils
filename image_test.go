package docker

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestPullImage(t *testing.T) {
	client := NewDockerClient()
	err := client.PullImage("docker.io/library/alpine")
	Nil(t, err)
}
