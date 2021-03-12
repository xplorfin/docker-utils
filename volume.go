package docker

import (
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
)

// RemoveContainer removes a stopped a container (and it's volumes) by id if it exists
// see Client.ContainerRemove or https://docs.docker.com/engine/reference/commandline/container_rm/
// for details
func (c *Client) RemoveContainer(id string) error {
	err := c.ContainerRemove(c.Ctx, id, types.ContainerRemoveOptions{
		RemoveVolumes: true,
	})
	return err
}

// StopContainer stops a container (but does not remove it) by id if it exists
// see Client.ContainerStop or https://docs.docker.com/engine/reference/commandline/container_stop/
// for details
func (c *Client) StopContainer(id string) error {
	timeout := time.Second * 5
	err := c.ContainerStop(c.Ctx, id, &timeout)
	return err
}

// RemoveVolume removes a volume by name if it exists
// see Client.VolumeRemove or https://docs.docker.com/engine/reference/commandline/volume_rm/
// for details
func (c *Client) RemoveVolume(name string) error {
	err := c.VolumeRemove(c.Ctx, name, true)
	return err
}

// VolumeExists checks if a volume exists by name
// see Client.VolumeInspect or https://docs.docker.com/engine/reference/commandline/volume_inspect/
// for details
func (c *Client) VolumeExists(name string) bool {
	_, err := c.VolumeInspect(c.Ctx, name)
	return err == nil
}

// CreateVolume creates a new docker volume with a unique name
// see Client.VolumeCreate or https://docs.docker.com/engine/extend/plugins_volume/
// for details
func (c *Client) CreateVolume(name string) (err error) {
	if c.VolumeExists(name) {
		err := c.RemoveVolume(name)
		if err != nil {
			return err
		}
	}
	vol, err := c.VolumeCreate(c.Ctx, volume.VolumeCreateBody{
		Driver: Driver,
		Name:   name,
		Labels: map[string]string{
			"sessionId": c.SessionID,
		},
	})
	if err != nil {
		return err
	}
	_ = vol

	return nil
}
