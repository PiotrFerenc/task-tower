package folder

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"os"
	"path/filepath"
)

type deleteFolder struct {
	config            *configuration.Config
	folderName        actions.Property
	deletedFolderPath actions.Property
}

func DeleteFolder(config *configuration.Config) actions.Action {
	return &deleteFolder{
		config: config,
		folderName: actions.Property{
			Name:        "folderName",
			Type:        actions.Text,
			Description: "The name of the folder to be deleted",
			DisplayName: "Folder Name",
			Validation:  "required",
		},
		deletedFolderPath: actions.Property{
			Name:        "createdFolderPath",
			Type:        actions.Text,
			Description: "The path where the new folder was created",
			DisplayName: "Created Folder Path",
			Validation:  "",
		},
	}
}

func (action *deleteFolder) GetCategoryName() string {
	return "folder"
}

func (action *deleteFolder) Inputs() []actions.Property {
	return []actions.Property{
		action.folderName,
	}
}

func (action *deleteFolder) Outputs() []actions.Property {
	return []actions.Property{}
}

// Execute performs the deletion of a folder with the provided folder name.
//
// Parameters:
//
//	message: The process message containing the necessary information.
//
// Returns:
//
//	types.Process: The updated process message.
//	error: An error if the deletion fails.
func (action *deleteFolder) Execute(message types.Process) (types.Process, error) {
	folderName, err := action.folderName.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}

	folderPath := filepath.Join(action.config.Folder.TmpFolder, folderName)

	err = os.RemoveAll(folderPath)
	if err != nil {
		return types.Process{}, err
	}
	message.SetString(action.deletedFolderPath.Name, folderPath)
	return message, nil
}
