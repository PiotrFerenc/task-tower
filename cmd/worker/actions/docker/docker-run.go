package docker

import (
	"context"
	"errors"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"strings"
)

func CreateDockerRun() actions.Action {
	return &dockerContainer{}
}

type dockerContainer struct {
}

const (
	imageName = "image"
	ports     = "ports"
)

func (c *dockerContainer) Inputs() []actions.Property {
	return []actions.Property{
		{
			Name: imageName,
			Type: "text",
		},
		{
			Name: ports,
			Type: "text",
		},
	}
}
func (c *dockerContainer) Outputs() []actions.Property {
	return []actions.Property{
		{Name: "container-id",
			Type: "text",
		},
	}
}

func (c *dockerContainer) Execute(process types.Pipeline) (types.Pipeline, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return process, err
	}
	imageName, err := process.GetString(imageName)

	if err != nil {
		return process, err
	}
	ports, err := process.GetString(ports)
	if err != nil {
		return process, err
	}

	portMap, portSet, err := mapPorts(ports)
	if err != nil {
		return process, err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:        imageName,
		ExposedPorts: portSet,
	}, &container.HostConfig{
		PortBindings: portMap,
	}, nil, nil, "")
	if err != nil {
		return process, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return process, err
	}

	process.SetString("container-id", resp.ID)
	return process, nil
}

func mapPorts(ports string) (nat.PortMap, nat.PortSet, error) {
	portPairs := strings.Split(ports, ", ")
	portMap := nat.PortMap{}
	portSet := nat.PortSet{}
	for _, portPair := range portPairs {
		portPair = strings.TrimSpace(portPair)
		if !strings.Contains(portPair, ":") {
			return nil, nil, errors.New("Invalid ports")
		}
		portParts := strings.Split(portPair, ":")
		containerPort := portParts[0]
		hostPort := portParts[1]
		portMap[nat.Port(containerPort)] = []nat.PortBinding{
			{
				HostIP:   "0.0.0.0",
				HostPort: hostPort,
			},
		}
		portSet[nat.Port(containerPort)] = struct{}{}
	}

	return portMap, portSet, nil
}
