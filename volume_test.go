package docker

import (
	"testing"

	. "github.com/stretchr/testify/assert"

	"github.com/brianvoe/gofakeit/v5"
)

func TestVolumeLifecycle(t *testing.T) {
	client := NewDockerClient()
	volName := gofakeit.LoremIpsumWord()
	err := client.CreateVolume(volName)
	Nil(t, err)
	// attempt to overwrite the same volume
	err = client.CreateVolume(volName)
	Nil(t, err)

	if !client.VolumeExists(volName) {
		t.Errorf("expected volume %s to exist", volName)
	}

	err = client.RemoveVolume(volName)
	Nil(t, err)
}
