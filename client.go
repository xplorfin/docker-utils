package docker

import (
	"context"

	"github.com/docker/docker/api/types"

	"github.com/google/uuid"

	"github.com/docker/docker/client"
)

// Driver defines the supported volume driver by docker-utils
// docker provides a number of different volume drivers for
// interacting with persistent disks. Since this library is
// oriented toward testing, we use local only here. For more
// info, see: https://docs.docker.com/engine/extend/plugins_volume/
const Driver = "local"

// Client contains context information for the current session
// to allow easy teardown at the end of the session
type Client struct {
	// pointer to the original docker client
	*client.Client
	// Ctx is the context object used for interacting with docker
	Ctx context.Context
	// SessionID is the identifier used for all resources created during
	// this session (as a v4 uuid)
	SessionID string
	// configured authentication goes here
	RegistryAuth *types.AuthConfig
}

// NewDockerClient creates a new Client and initializes it with a sessionId.
// Note: in some cases of failure, deferred methods like Client.TeardownSession
// may not be called. It may be useful to print out the session id and tell the
// user how to remove all items associated with a session in this scenario
func NewDockerClient() Client {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	sessionID, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return Client{
		Client:    cli,
		Ctx:       ctx,
		SessionID: sessionID.String(),
	}
}
