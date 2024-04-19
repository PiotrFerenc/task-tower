package folder

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"os"
	"path/filepath"
)

type createFolder struct {
	config            *configuration.Config
	folderName        actions.Property
	createdFolderPath actions.Property
}

func CreateFolder(config *configuration.Config) actions.Action {
	return &createFolder{
		config: config,
		folderName: actions.Property{
			Name:        "folderName",
			Type:        actions.Text,
			Description: "The name of the folder",
			DisplayName: "Folder Name",
			Validation:  "required",
		},
		createdFolderPath: actions.Property{
			Name:        "createdFolderPath",
			Type:        actions.Text,
			Description: "The path where the new folder was created",
			DisplayName: "Created Folder Path",
			Validation:  "",
		},
	}
}

func (action *createFolder) GetCategoryName() string {
	return "folder"
}

func (action *createFolder) Inputs() []actions.Property {
	return []actions.Property{
		action.folderName,
	}
}

func (action *createFolder) Outputs() []actions.Property {
	return []actions.Property{
		action.createdFolderPath,
	}
}

// Execute executes the createFolder action by creating a folder with the provided name in the temporary folder.
// It returns the updated process message and error, if any.
//
// Parameters:
//
//	message: The process message to operate on.
//
// Returns:
//
//	types.Process: The updated process message.
//	error:         An error if any occurred during execution.
func (action *createFolder) Execute(message types.Process) (types.Process, error) {
	folderName, err := action.folderName.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}

	folderPath := filepath.Join(action.config.Folder.TmpFolder, folderName)
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return types.Process{}, err
	}
	message.SetString(action.createdFolderPath.Name, folderPath)
	return message, nil
}
