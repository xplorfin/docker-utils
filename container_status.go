package docker

// ContainerStatus defines constants for container statuses in docker
// the standard docker library requires us to make string comparisons that are wrapped here
// to make comparators easier to use.
// see statuses here: https://docs.docker.com/engine/reference/commandline/ps/
type ContainerStatus string

// container is running
const (
	// ContainerCreated indicates a container has been created
	// this is the status after container create
	// see: https://docs.docker.com/engine/reference/commandline/container_create/ for details
	ContainerCreated ContainerStatus = "created"
	// ContainerRunning indicates a command inside a container is running
	// see: https://docs.docker.com/engine/reference/commandline/container_run/
	ContainerRunning ContainerStatus = "running"
	// ContainerPaused indicates a container is paused and is not available
	// see: https://docs.docker.com/engine/reference/commandline/container_pause/
	ContainerPaused ContainerStatus = "paused"
	// ContainerRestarting indicates a container is in the process of restarting
	// and is not available
	// see: https://docs.docker.com/engine/reference/commandline/container_restart/
	ContainerRestarting ContainerStatus = "restarting"
	// ContainerExited indicates the command run in the container has exited
	ContainerExited ContainerStatus = "exited"
	// ContainerDead is a container that cannot be restarted
	// see: https://git.io/JqV6W
	ContainerDead ContainerStatus = "dead"
	// ContainerRemoving is a container that is in the process of being removed
	// see: https://docs.docker.com/engine/reference/commandline/rm/
	ContainerRemoving ContainerStatus = "removing"
	// ContainerUnknown is a container status that is unknown.
	ContainerUnknown ContainerStatus = "unknown"
)
