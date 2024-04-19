package dotnet

import (
	"context"
	"fmt"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/Container"
	"github.com/PiotrFerenc/mash2/internal/types"
)

type solutionFile struct {
	ProjectName actions.Property
	ContainerId actions.Property
}

func CreateDotnetSolutionAction() actions.Action {
	return &solutionFile{
		ProjectName: actions.Property{
			Name:        "SolutionName",
			Type:        "string",
			Description: "Name of the solution.",
			DisplayName: "Solution Name",
			Validation:  "Required",
		},
		ContainerId: actions.Property{
			Name:        "ContainerId",
			Type:        "string",
			Description: "Unique identifier of the container.",
			DisplayName: "Container ID",
			Validation:  "Required",
		},
	}
}

func (cl *solutionFile) Inputs() []actions.Property {
	return []actions.Property{
		cl.ProjectName,
	}
}
func (cl *solutionFile) Outputs() []actions.Property {
	return []actions.Property{cl.ContainerId}
}
func (cl *solutionFile) GetCategoryName() string {
	return "dotnet"
}

// Execute executes a solutionFile process by starting a container with the specified image and environment variables.
//
// Parameters:
//
//	process: The process to be executed.
//
// Returns:
//
//	types.Process: The modified process object.
//	error: An error if any occurred during the execution.
func (cl *solutionFile) Execute(process types.Process) (types.Process, error) {
	ctx := context.Background()
	projectName, err := cl.ProjectName.GetStringFrom(&process)
	if err != nil {
		return process, err
	}
	env := fmt.Sprintf("SOLUTION_NAME=%s", projectName)
	imageName := "dotnet-solutionFile"
	vol := fmt.Sprintf("/dashboard/appuser/%s:/data", projectName)
	containerId, err := Container.StartContainer(imageName, []string{env}, []string{vol}, ctx)

	process.SetString(cl.ContainerId.Name, containerId)
	return process, nil
}
