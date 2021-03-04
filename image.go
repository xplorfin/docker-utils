package docker

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"os"

	"github.com/docker/docker/api/types"
)

func (c Client) generatePullConfig(image string) (options types.ImagePullOptions) {
	if c.RegistryAuth != nil {
		encodedJson, err := json.Marshal(c.RegistryAuth)
		if err != nil {
			panic(err)
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJson)
		return types.ImagePullOptions{
			RegistryAuth: authStr,
		}
	}
	return options
}

// pull a docker image
func (c Client) PullImage(image string) error {
	reader, err := c.ImagePull(c.Ctx, image, c.generatePullConfig(image))
	if err != nil {
		return err
	}
	_, err = io.Copy(os.Stdout, reader)
	return err
}
