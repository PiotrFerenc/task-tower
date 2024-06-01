package docker

import (
	"context"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/Container"
	"github.com/PiotrFerenc/mash2/internal/types"
)

type dockerRemoveContainer struct {
	containerId actions.Property
}

func CreateDockerRemove() actions.Action {
	return &dockerRemoveContainer{
		containerId: actions.Property{
			Name:        "id",
			Type:        actions.Text,
			Description: "The unique identifier for the Docker container to be removed",
			DisplayName: "Container ID",
			Validation:  "required",
		},
	}
}

func (d *dockerRemoveContainer) GetCategoryName() string {
	return "docker"
}

func (d *dockerRemoveContainer) Inputs() []actions.Property {
	return []actions.Property{
		d.containerId,
	}
}

func (d *dockerRemoveContainer) Outputs() []actions.Property {
	return []actions.Property{}
}

func (d *dockerRemoveContainer) Execute(process types.Process) (types.Process, error) {
	ctx := context.Background()

	containerId, err := d.containerId.GetStringFrom(&process)
	if err != nil {
		return process, err
	}

	err = Container.RemoveContainer(containerId, ctx)
	if err != nil {
		return process, err
	}

	return process, nil
}
