package docker

import "github.com/docker/docker/api/types"

// CreateNetwork creates a docker network that can
// be used for communication in between containers
// see: https://docs.docker.com/network/ for details
func (c Client) CreateNetwork(name string) (networkID string, err error) {
	networkResponse, err := c.NetworkCreate(c.Ctx, name,
		types.NetworkCreate{
			CheckDuplicate: true,
			Attachable:     true,
			Labels:         c.getSessionLabels(),
		})
	if err != nil {
		return "", err
	}
	return networkResponse.ID, err
}

// RemoveNetwork removes a network by id
func (c Client) RemoveNetwork(networkID string) error {
	return c.NetworkRemove(c.Ctx, networkID)
}