package docker

import (
	"context"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/Container"
	"github.com/PiotrFerenc/mash2/internal/types"
	"strings"
)

func CreateDockerRun() actions.Action {
	return &dockerContainer{
		imageName: actions.Property{
			Name:        "image",
			Type:        actions.Text,
			Description: "The name of the Docker image to use",
			DisplayName: "Image Name",
			Validation:  "required",
		},
		env: actions.Property{
			Name:        "e",
			Type:        actions.Text,
			Description: "The environment variables for the Docker container",
			DisplayName: "Environment Variables",
			Validation:  "required",
		},
		vol: actions.Property{
			Name:        "v",
			Type:        actions.Text,
			Description: "Volume",
			DisplayName: "Environment Variables",
			Validation:  "required",
		},
		containerId: actions.Property{
			Name:        "id",
			Type:        actions.Text,
			Description: "The unique identifier for the Docker container",
			DisplayName: "Container ID",
			Validation:  "",
		},
	}
}

type dockerContainer struct {
	imageName   actions.Property
	containerId actions.Property
	env         actions.Property
	vol         actions.Property
}

func (d *dockerContainer) GetCategoryName() string {
	return "docker"
}

func (d *dockerContainer) Inputs() []actions.Property {
	return []actions.Property{
		d.imageName, d.env, d.vol,
	}
}
func (d *dockerContainer) Outputs() []actions.Property {
	return []actions.Property{
		d.containerId,
	}
}

// Execute executes a Docker container.
//
// Parameters:
//
//	process: The process that contains the required parameters for executing the container.
//
// Returns:
//
//	types.Process: The updated process with the container ID set.
//	error: An error if the container execution fails.
func (d *dockerContainer) Execute(process types.Process) (types.Process, error) {
	ctx := context.Background()
	imageName, err := d.imageName.GetStringFrom(&process)
	if err != nil {
		return process, err
	}
	vol, err := d.vol.GetStringFrom(&process)
	if err != nil {
		return process, err
	}

	envParameters, err := d.env.GetStringFrom(&process)
	if err != nil {
		return process, err
	}
	env := strings.Split(envParameters, ",")
	containerId, err := Container.StartContainer(imageName, env, []string{vol}, ctx)

	process.SetString(d.containerId.Name, containerId)
	return process, nil
}
