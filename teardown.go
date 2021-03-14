package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/hashicorp/go-multierror"
)

// TeardownSession removes all resources created during the session
func (c *Client) TeardownSession() (err error) {
	err = c.TeardownSessionContainers()
	if err != nil {
		err = multierror.Append(err)
	}
	err = c.TeardownSessionVolumes()
	if err != nil {
		err = multierror.Append(err)
	}
	err = c.TeardownSessionNetworks()
	if err != nil {
		err = multierror.Append(err)
	}
	return err
}

// TeardownSessionVolumes removes containers created in the current session using RemoveVolume
// the current session is determined based on the Client so this method should be called
// once per Client
func (c *Client) TeardownSessionVolumes() (err error) {
	volumes, err := c.VolumeList(c.Ctx, filterByLabels(c.getSessionLabels()))

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

// TeardownSessionNetworks removes containers created in the current session using StopContainer
// and RemoveContainer the current session is determined based on the Client
// so this method should be called once per Client
func (c *Client) TeardownSessionNetworks() (err error) {
	networks, err := c.NetworkList(c.Ctx, types.NetworkListOptions{Filters: filterByLabels(c.getSessionLabels())})
	if err != nil {
		return err
	}

	for _, network := range networks {
		err = c.RemoveNetwork(network.ID)
		if err != nil {
			return err
		}
	}
	return err
}

// TeardownSessionContainers removes containers created in the current session using StopContainer
// and RemoveContainer the current session is determined based on the Client
// so this method should be called once per Client
func (c *Client) TeardownSessionContainers() (err error) {
	filter := types.ContainerListOptions{
		All:     true,
		Filters: filterByLabels(c.getSessionLabels()),
	}

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
