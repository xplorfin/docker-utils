package docker

import (
	. "github.com/stretchr/testify/assert"
	"testing"
)

func TestPullImage(t *testing.T) {
	client := NewDockerClient()
	err := client.PullImage("docker.io/library/alpine")
	Nil(t, err)
}
