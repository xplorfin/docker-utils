package docker

import "testing"

func TestNewClient(t *testing.T) {
	client := NewDockerClient()
	_, err := client.Info(client.Ctx)
	if err != nil {
		t.Error(err)
	}
}
