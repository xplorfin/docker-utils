package docker

import (
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// removes *everything* created during this session
func (c *Client) TeardownSession() (err error) {
	err = c.TeardownSessionContainers()
	if err != nil {
		return err
	}
	return c.TeardownSessionVolumes()
}

// remove volumes created in this session
func (c *Client) TeardownSessionVolumes() error {
	volumes, err := c.VolumeList(c.Ctx, filters.NewArgs(filters.KeyValuePair{
		Key:   "label",
		Value: fmt.Sprintf("sessionId=%s", c.SessionId),
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

// remove containers created in this session
func (c *Client) TeardownSessionContainers() error {
	filter := types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: fmt.Sprintf("sessionId=%s", c.SessionId),
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
