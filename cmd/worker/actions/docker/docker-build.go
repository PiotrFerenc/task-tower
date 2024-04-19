package docker

import (
	"context"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/Container"
	"github.com/PiotrFerenc/mash2/internal/types"
)

func CreateDockerBuild() actions.Action {
	return &dockerImage{
		dockerfile: actions.Property{
			Name:        "dockerfile",
			Type:        actions.Text,
			Description: "The dockerfile to use for building the image",
			DisplayName: "Dockerfile",
			Validation:  "required",
		},
		tags: actions.Property{
			Name:        "tags",
			Type:        actions.Text,
			Description: "The tags to apply to the built image",
			DisplayName: "Tags",
			Validation:  "",
		},
		imageId: actions.Property{
			Name:        "id",
			Type:        actions.Text,
			Description: "The unique identifier for the built Docker image",
			DisplayName: "Image ID",
			Validation:  "",
		},
	}
}

type dockerImage struct {
	dockerfile actions.Property
	tags       actions.Property
	imageId    actions.Property
}

func (d *dockerImage) GetCategoryName() string {
	return "docker"
}

func (d *dockerImage) Inputs() []actions.Property {
	return []actions.Property{
		d.dockerfile, d.tags,
	}
}

func (d *dockerImage) Outputs() []actions.Property {
	return []actions.Property{
		d.imageId,
	}
}

// Execute builds a Docker image using the provided Dockerfile content and tags.
//
// Parameters:
//
//	process: The process object containing the Dockerfile content and tags.
//
// Returns:
//
//	types.Process: The updated process object with the image ID set.
//	error: An error if the image build process fails.
func (d *dockerImage) Execute(process types.Process) (types.Process, error) {
	ctx := context.Background()
	dockerfile, err := d.dockerfile.GetStringFrom(&process)
	if err != nil {
		return process, err
	}
	tags, err := d.tags.GetStringFrom(&process)
	if err != nil {
		return process, err
	}
	imageId, err := Container.BuildImage(dockerfile, tags, ctx)
	process.SetString(d.imageId.Name, imageId)
	return process, nil
}
