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

type ExecResult struct {
	// stdout output
	StdOut string
	// stderr output
	StdErr string
	// exit code of the process: https://tldp.org/LDP/abs/html/exitcodes.html
	ExitCode int
}

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

// exec method that allows you to pass your own config
// make sure you attach stderr and stdout if you need them
func (c *Client) ExecRaw(containerId string, config types.ExecConfig) (ExecResult, error) {
	execution, err := c.ContainerExecCreate(c.Ctx, containerId, config)
	if err != nil {
		return ExecResult{}, err
	}

	return c.InspectExecResp(c.Ctx, execution.ID)
}

func (c *Client) Exec(containerId string, cmd []string) (ExecResult, error) {
	config := types.ExecConfig{
		AttachStderr: true,
		AttachStdout: true,
		Cmd:          cmd,
	}

	return c.ExecRaw(containerId, config)
}

// wrapper around container create. Pulls image first/adds context
func (c Client) CreateContainer(config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *specs.Platform, containerName string) (container.ContainerCreateCreatedBody, error) {
	err := c.PullImage(config.Image)
	if err != nil {
		return container.ContainerCreateCreatedBody{}, err
	}
	if config.Labels == nil {
		config.Labels = make(map[string]string)
	}
	config.Labels["sessionId"] = c.SessionId
	return c.ContainerCreate(c.Ctx, config, hostConfig, networkingConfig, platform, containerName)
}

// copy container lgos to stdout
func (c Client) PrintContainerLogs(containerId string) {
	go func() {
		statusCh, errCh := c.ContainerWait(c.Ctx, containerId, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}

		out, err := c.ContainerLogs(c.Ctx, containerId, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			// container is dead
			return
		}
		_, _ = stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}()

}
