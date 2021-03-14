package docker

import (
	"bytes"
	"context"
	"io/ioutil"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

// ExecResult gets the result of a executed docker command
type ExecResult struct {
	// StdOut displays what is send to stdout by the container. This is normally in docker logs
	// See: https://docs.docker.com/config/containers/logging/
	StdOut string
	// StdOut displays what is send to stderr by the container. This is normally in docker logs
	// See: https://docs.docker.com/config/containers/logging/
	StdErr string
	// ExitCode of the process. See https://tldp.org/LDP/abs/html/exitcodes.html for details
	ExitCode int
}

// InspectExecResp copies container execution results into an ExecResult
// object after thje command is finished executing. Returns error on failure
// See: https://docs.docker.com/engine/reference/commandline/attach/ for details
func (c *Client) InspectExecResp(ctx context.Context, id string) (ExecResult, error) {
	var execResult ExecResult

	resp, err := c.ContainerExecAttach(ctx, id, types.ExecStartCheck{})
	if err != nil {
		return execResult, err
	}
	defer resp.Close()

	// read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, resp.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			return execResult, err
		}
		break

	case <-ctx.Done():
		return execResult, ctx.Err()
	}

	stdout, err := ioutil.ReadAll(&outBuf)
	if err != nil {
		return execResult, err
	}
	stderr, err := ioutil.ReadAll(&errBuf)
	if err != nil {
		return execResult, err
	}

	res, err := c.ContainerExecInspect(ctx, id)
	if err != nil {
		return execResult, err
	}

	execResult.ExitCode = res.ExitCode
	execResult.StdOut = string(stdout)
	execResult.StdErr = string(stderr)
	return execResult, nil
}

// ExecRaw allows you to pass a custom types.ExecConfig to a container
// running your command. This method will then attach to the container and return ExecResult
// using InspectExecResp. See https://docs.docker.com/engine/reference/commandline/exec/
// or Client.ContainerExecCreate for details
func (c *Client) ExecRaw(containerID string, config types.ExecConfig) (ExecResult, error) {
	execution, err := c.ContainerExecCreate(c.Ctx, containerID, config)
	if err != nil {
		return ExecResult{}, err
	}

	return c.InspectExecResp(c.Ctx, execution.ID)
}

// Exec executes a command cmd in a container using ExecRaw
// This method will then attach to the container and return ExecResult
// using InspectExecResp. See https://docs.docker.com/engine/reference/commandline/exec/ for details
func (c *Client) Exec(containerID string, cmd []string) (ExecResult, error) {
	config := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}

	return c.ExecRaw(containerID, config)
}

// CreateContainer creates a container with a given *container.Config
// making sure to pull the image using PullImage if not present and session the SessionId
// for variadic teardown. See https://docs.docker.com/engine/reference/commandline/create/ for details
func (c Client) CreateContainer(config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	err := c.PullImage(config.Image)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}
	if config.Labels == nil {
		config.Labels = make(map[string]string)
	}
	config.Labels["sessionId"] = c.SessionID
	return c.ContainerCreate(c.Ctx, config, hostConfig, networkingConfig, platform, containerName)
}

// PrintContainerLogs copies logs for a running, non-executing command to stdout
// See https://docs.docker.com/engine/reference/commandline/logs/ for details
// TODO: in the futre this should be more flexible
func (c Client) PrintContainerLogs(containerID string) {
	go func() {
		statusCh, errCh := c.ContainerWait(c.Ctx, containerID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}

		out, err := c.ContainerLogs(c.Ctx, containerID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			// container is dead
			return
		}
		_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}()

}

// ContainerStatus fetches a container's status and normalizes it
// see: https://docs.docker.com/engine/reference/commandline/inspect/
func (c Client) ContainerStatus(containerID string) (ContainerStatus, error) {
	inspectedContainer, err := c.ContainerInspect(c.Ctx, containerID)
	if err != nil {
		return "", err
	}

	return ContainerStatus(inspectedContainer.ContainerJSONBase.State.Status), err
}
