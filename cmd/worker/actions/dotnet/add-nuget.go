package dotnet

import (
	"context"
	"fmt"
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/Container"
	"github.com/PiotrFerenc/mash2/internal/types"
)

func AddPackageToProject() actions.Action {
	return &addToProject{
		projectPath: actions.Property{
			Name:        "projectPath",
			Type:        actions.Text,
			Description: "The path to the project where the package needs to be added",
			DisplayName: "Project Path",
			Validation:  "required",
		},
		packageName: actions.Property{
			Name:        "packageName",
			Type:        actions.Text,
			Description: "The name of the package that needs to be added to the project",
			DisplayName: "Package Name",
			Validation:  "required",
		},
		containerId: actions.Property{
			Name:        "ContainerId",
			Type:        "string",
			Description: "Unique identifier of the container.",
			DisplayName: "Container ID",
			Validation:  "Required",
		},
	}
}

type addToProject struct {
	projectPath actions.Property
	packageName actions.Property
	containerId actions.Property
}

func (a addToProject) Inputs() []actions.Property {
	return []actions.Property{
		a.projectPath, a.packageName,
	}
}

func (a addToProject) Outputs() []actions.Property {
	return []actions.Property{
		a.containerId,
	}
}

func (a addToProject) GetCategoryName() string {
	return "docker"
}

func (a addToProject) Execute(process types.Process) (types.Process, error) {
	ctx := context.Background()

	projectPath, err := a.projectPath.GetStringFrom(&process)
	if err != nil {
		return process, err
	}
	packageName, err := a.packageName.GetStringFrom(&process)
	if err != nil {
		return process, err
	}

	envProjectPath := fmt.Sprintf("PROJECT_PATH=%s", projectPath)
	envPackageName := fmt.Sprintf("NUGET_NAME=%s", packageName)
	imageName := "add-nuget"
	vol := fmt.Sprintf("/app/%s:/data", packageName)

	containerId, err := Container.StartContainer(imageName, []string{envProjectPath, envPackageName}, []string{vol}, ctx)

	process.SetString(a.containerId.Name, containerId)
	return process, nil
}
