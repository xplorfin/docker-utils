package docker

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/docker/docker/api/types"
)

// generatePullConfig generates a pull configuration with the RegistryAuth if present
// see https://docs.docker.com/engine/reference/commandline/login/ for details
func (c Client) generatePullConfig() (options types.ImagePullOptions) {
	if c.RegistryAuth != nil {
		encodedJSON, err := json.Marshal(c.RegistryAuth)
		if err != nil {
			panic(err)
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		return types.ImagePullOptions{
			RegistryAuth: authStr,
		}
	}
	return options
}

// PullImage pulls a docker image by name. If RegistryAuth is specified, it is used here
// see: https://docs.docker.com/engine/reference/commandline/pull/ for details
func (c Client) PullImage(image string) error {
	reader, err := c.ImagePull(c.Ctx, image, c.generatePullConfig())
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, reader)
	return err
}
