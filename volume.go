package docker

import (
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
)

func (c *Client) RemoveContainer(id string) error {
	err := c.ContainerRemove(c.Ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
	})
	return err
}

func (c *Client) StopContainer(id string) error {
	timeout := time.Second * 5
	err := c.ContainerStop(c.Ctx, id, &timeout)
	return err
}

func (c *Client) RemoveVolume(name string) error {
	err := c.VolumeRemove(c.Ctx, name, true)
	return err
}

// check if a volume exists
func (c *Client) VolumeExists(name string) bool {
	_, err := c.VolumeInspect(c.Ctx, name)
	return err == nil
}

// create a new volume
func (c *Client) CreateVolume(name string) error {
	if c.VolumeExists(name) {
		c.RemoveVolume(name)
	}
	vol, err := c.VolumeCreate(c.Ctx, volume.VolumeCreateBody{
		Driver: Driver,
		Name:   name,
		Labels: map[string]string{
			"sessionId": c.SessionId,
		},
	})
	if err != nil {
		return err
	}
	_ = vol

	return nil
}
