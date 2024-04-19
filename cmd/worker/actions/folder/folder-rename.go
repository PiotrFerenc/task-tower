package folder

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"os"
	"path/filepath"
)

type renameFolder struct {
	config       *configuration.Config
	oldFolder    actions.Property
	newFolder    actions.Property
	pathToFolder actions.Property
}

func RenameFolder(config *configuration.Config) actions.Action {
	return &renameFolder{
		config: config,
		oldFolder: actions.Property{
			Name:        "oldfolderName",
			Type:        actions.Text,
			Description: "The name of the folder to be renamed",
			DisplayName: "Original Folder Name",
			Validation:  "required",
		},
		newFolder: actions.Property{
			Name:        "newfolderName",
			Type:        actions.Text,
			Description: "The new name for the folder",
			DisplayName: "Target Folder Name",
			Validation:  "required",
		},
		pathToFolder: actions.Property{
			Name:        "renamedFolderPath",
			Type:        actions.Text,
			Description: "Path to the renamed folder",
			DisplayName: "Renamed Folder Path",
			Validation:  "",
		},
	}
}

func (action *renameFolder) GetCategoryName() string {
	return "folder"
}

func (action *renameFolder) Inputs() []actions.Property {
	return []actions.Property{
		action.oldFolder,
		action.newFolder,
	}
}

func (action *renameFolder) Outputs() []actions.Property {
	return []actions.Property{
		action.pathToFolder,
	}
}

// Execute renames a folder by moving it to a new location with a new name.
//
// Parameters:
//
//	message: The original Process object.
//
// Returns:
//
//	types.Process: The modified Process object with the updated folder path.
//	error: An error if the folder cannot be renamed.
func (action *renameFolder) Execute(message types.Process) (types.Process, error) {
	oldFolder, err := action.oldFolder.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}

	newFolder, err := action.newFolder.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}

	oldPath := filepath.Join(action.config.Folder.TmpFolder, oldFolder)
	newPath := filepath.Join(action.config.Folder.TmpFolder, newFolder)
	err = os.Rename(oldPath, newPath)
	if err != nil {
		return types.Process{}, err
	}

	message.SetString(action.pathToFolder.Name, newPath)
	return message, nil
}
