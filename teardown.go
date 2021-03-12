package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// TeardownSession removes all resources created during the session
func (c *Client) TeardownSession() (err error) {
	err = c.TeardownSessionContainers()
	if err != nil {
		return err
	}
	return c.TeardownSessionVolumes()
}

// TeardownSessionVolumes removes containers created in the current session using RemoveVolume
// the current session is determined based on the Client so this method should be called
// once per Client
func (c *Client) TeardownSessionVolumes() error {
	volumes, err := c.VolumeList(c.Ctx, filters.NewArgs(filters.KeyValuePair{
		Key:   "label",
		Value: fmt.Sprintf("sessionId=%s", c.SessionID),
	}))

	if err != nil {
		return err
	}

	for _, selectedVolume := range volumes.Volumes {
		err = c.RemoveVolume(selectedVolume.Name)
		if err != nil {
			return err
		}
	}

	return err
}

// TeardownSessionContainers removes containers created in the current session using StopContainer
// and RemoveContainer the current session is determined based on the Client
// so this method should be called once per Client
func (c *Client) TeardownSessionContainers() error {
	filter := types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: fmt.Sprintf("sessionId=%s", c.SessionID),
		}),
	}

	_ = filter
	containers, err := c.ContainerList(c.Ctx, filter)

	if err != nil {
		return err
	}

	for _, selectedContainer := range containers {
		err = c.StopContainer(selectedContainer.ID)
		if err != nil {
			return err
		}
		err = c.RemoveContainer(selectedContainer.ID)
		if err != nil {
			return err
		}
	}

	return err
}
