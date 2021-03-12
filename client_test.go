package docker

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewDockerClient()
	_, err := client.Info(client.Ctx)
	Nil(t, err)
}
