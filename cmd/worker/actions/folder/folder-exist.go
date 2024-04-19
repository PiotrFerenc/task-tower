package folder

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"os"
	"path/filepath"
)

type checkFolder struct {
	config     *configuration.Config
	folderName actions.Property
}

func CheckFolder(config *configuration.Config) actions.Action {
	return &checkFolder{
		config: config,
		folderName: actions.Property{
			Name:        "folderName",
			Type:        actions.Text,
			Description: "The name of the folder to be checked",
			DisplayName: "Folder Name",
			Validation:  "required",
		},
	}
}

func (action *checkFolder) GetCategoryName() string {
	return "folder"
}

func (action *checkFolder) Inputs() []actions.Property {
	return []actions.Property{
		action.folderName,
	}
}

func (action *checkFolder) Outputs() []actions.Property {
	return []actions.Property{}
}

// Execute executes the checkFolder action by checking if a folder with the provided folderName exists in the temporary folder.
//
// Parameters:
//
//	message: The input Process object containing the necessary data for executing the action.
//
// Returns:
//
//	Process: The updated Process object.
//	error:   Error if any occurred during the execution of the action.
func (action *checkFolder) Execute(message types.Process) (types.Process, error) {
	folderName, err := action.folderName.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}
	folderPath := filepath.Join(action.config.Folder.TmpFolder, folderName)
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return types.Process{}, err
	}

	return message, nil
}
