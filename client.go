package docker

import (
	"context"

	"github.com/docker/docker/api/types"

	"github.com/google/uuid"

	"github.com/docker/docker/client"
)

const Driver = "local"

type Client struct {
	*client.Client
	Ctx       context.Context
	SessionId string
	// configured authentication goes here
	RegistryAuth *types.AuthConfig
}

func NewDockerClient() Client {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	sessionId, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return Client{
		Client:    cli,
		Ctx:       ctx,
		SessionId: sessionId.String(),
	}
}
