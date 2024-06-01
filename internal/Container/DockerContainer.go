package Container

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"os"
)

// StartContainer starts a Docker container with the specified image name, environment variables,
// volume bindings, and context.
//
// Parameters:
//
//	imageName: The name of the Docker image to use for the container.
//	env: The environment variables to set for the container.
//	vol: The volume bindings for the container.
//	ctx: The context to use for the Docker client.
//
// Returns:
//
//	string: The ID of the container that was started.
//	error: An error if the container failed to start.
func StartContainer(imageName string, env, vol []string, ctx context.Context) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   env,
	}, &container.HostConfig{
		Binds: vol,
	}, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return "", err
	}
	return resp.ID, nil
}

// BuildImage builds a Docker image using the specified Dockerfile path, image name, and context.
//
// Parameters:
//
//	dockerfilePath: The path to the Dockerfile used for building the image.
//	imageName:      The name of the image to be built.
//	ctx:            The context to use for the Docker client.
//
// Returns:
//
//	string: The ID of the built image.
//	error:  An error if the image build process fails.
func BuildImage(dockerfilePath, imageName string, ctx context.Context) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return "", err
	}

	buildContext, err := os.Open(dockerfilePath)
	if err != nil {
		return "", err
	}
	defer buildContext.Close()

	options := types.ImageBuildOptions{
		Dockerfile: dockerfilePath,
		Tags:       []string{imageName},
		Remove:     true,
	}

	response, err := cli.ImageBuild(ctx, buildContext, options)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	//body, err := io.ReadAll(response.Body)
	//if err != nil {
	//	return "", err
	//}
	imgID, _, err := cli.ImageInspectWithRaw(ctx, imageName)
	if err != nil {
		return "", err
	}

	return imgID.ID, nil
}

// RemoveContainer removes a Docker container with the specified container ID and context.
//
// Parameters:
//
//	containerId: The ID of the container to be removed.
//	ctx: The context to use for the Docker client.
//
// Returns:
//
//	error: An error if the container failed to be removed.
func RemoveContainer(containerId string, ctx context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	err = cli.ContainerRemove(ctx, containerId, container.RemoveOptions{
		RemoveVolumes: true,
		RemoveLinks:   true,
		Force:         true,
	})

	if err != nil {
		return err
	}

	return nil
}
