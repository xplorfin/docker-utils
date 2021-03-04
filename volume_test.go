package docker

import (
	"testing"

	"github.com/brianvoe/gofakeit/v5"
)

func TestVolumeLifecycle(t *testing.T) {
	client := NewDockerClient()
	volName := gofakeit.LoremIpsumWord()
	err := client.CreateVolume(volName)
	if err != nil {
		t.Error(err)
	}
	if !client.VolumeExists(volName) {
		t.Errorf("expected volume %s to exist", volName)
	}

	err = client.RemoveVolume(volName)
	if err != nil {
		t.Error(err)
	}

}
