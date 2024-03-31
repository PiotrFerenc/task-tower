package docker

import (
	"context"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

func CreateDockerRun() actions.Action {
	return &dockerContainer{}
}

type dockerContainer struct {
}

func (c *dockerContainer) Inputs() []actions.Property {
	output := make([]actions.Property, 1)
	output[0] = actions.Property{
		Name: "image",
		Type: "text",
	}
	return output
}
func (c *dockerContainer) Outputs() []actions.Property {
	output := make([]actions.Property, 1)
	output[0] = actions.Property{
		Name: "docker-container-id",
		Type: "text",
	}
	return output
}
func (c *dockerContainer) Execute(process types.Pipeline) (types.Pipeline, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return process, err
	}
	imageName, err := process.GetString("image")
	if err != nil {
		return process, err
	}
	out, err := cli.ImagePull(ctx, imageName, image.PullOptions{})
	if err != nil {
		panic(err)
	}
	defer out.Close()

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
	}, nil, nil, nil, "")
	if err != nil {
		return process, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return process, err
	}

	out, err = cli.ContainerLogs(ctx, resp.ID, container.LogsOptions{ShowStdout: true})
	if err != nil {
		return process, err
	}
	defer out.Close()

	process.SetString("docker-container-id", resp.ID)
	return process, nil
}
