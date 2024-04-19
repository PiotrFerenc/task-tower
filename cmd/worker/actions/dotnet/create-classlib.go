package dotnet

import (
	"context"
	"fmt"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/Container"
	"github.com/PiotrFerenc/mash2/internal/types"
)

type classlib struct {
	ProjectName actions.Property
	Containerid actions.Property
}

func CreateDotnetClassLibAction() actions.Action {
	return &classlib{
		ProjectName: actions.Property{
			Name:        "ProjectName",
			Type:        "string",
			Description: "Name of the project.",
			DisplayName: "Project Name",
			Validation:  "Required",
		},
		Containerid: actions.Property{
			Name:        "ContainerId",
			Type:        "string",
			Description: "Unique identifier of the container.",
			DisplayName: "Container ID",
			Validation:  "Required",
		},
	}
}

func (cl *classlib) Inputs() []actions.Property {
	return []actions.Property{
		cl.ProjectName,
	}
}
func (cl *classlib) Outputs() []actions.Property {
	return []actions.Property{}
}
func (cl *classlib) GetCategoryName() string {
	return "dotnet"
}

// Execute executes a classlib process by starting a container with the specified image and environment variables.
//
// Parameters:
//
//	process: The process to be executed.
//
// Returns:
//
//	types.Process: The modified process object.
//	error: An error if any occurred during the execution.
func (cl *classlib) Execute(process types.Process) (types.Process, error) {
	ctx := context.Background()
	projectName, err := cl.ProjectName.GetStringFrom(&process)
	if err != nil {
		return process, err
	}
	env := fmt.Sprintf("PROJECT_NAME=%s", projectName)
	imageName := "dotnet-classlib"
	vol := fmt.Sprintf("/dashboard/appuser/%s:/data", projectName)
	containerId, err := Container.StartContainer(imageName, []string{env}, []string{vol}, ctx)

	process.SetString(cl.Containerid.Name, containerId)
	return process, nil
}
