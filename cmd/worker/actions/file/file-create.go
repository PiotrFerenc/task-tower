package file

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"os"
	"path/filepath"
)

type contentToFile struct {
	config          *configuration.Config
	fileName        actions.Property
	content         actions.Property
	createdFilePath actions.Property
}

func CreateContentToFile(config *configuration.Config) actions.Action {
	return &contentToFile{
		config: config,
		fileName: actions.Property{
			Name:        "fileName",
			Type:        actions.Text,
			Description: "The name of the file",
			DisplayName: "File Name",
			Validation:  "required",
		},
		content: actions.Property{
			Name:        "content",
			Type:        actions.Text,
			Description: "The content to be written to the file",
			DisplayName: "File Content",
			Validation:  "required",
		},
		createdFilePath: actions.Property{
			Name:        "createdFilePath",
			Type:        actions.Text,
			Description: "The path where the new file was created",
			DisplayName: "Created File Path",
			Validation:  "",
		},
	}
}

func (action *contentToFile) GetCategoryName() string {
	return "file"
}
func (action *contentToFile) Inputs() []actions.Property {
	return []actions.Property{
		action.fileName, action.content,
	}
}

func (action *contentToFile) Outputs() []actions.Property {
	return []actions.Property{
		action.createdFilePath,
	}
}

// Execute executes the contentToFile action by writing the content to a file with the provided file name.
// The file is created in the TmpFolder specified in the configuration.
// The file path is stored in the createdFilePath property of the message.
//
// Parameters:
//
//	message: The process message containing the necessary data for the action to execute.
//
// Returns:
//
//	types.Process: The updated process message after executing the action.
//	error: An error if the action encounters any issues during execution.
func (action *contentToFile) Execute(message types.Process) (types.Process, error) {
	fileName, err := action.fileName.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}

	content, err := action.content.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}

	filePath := filepath.Join(action.config.Folder.TmpFolder, fileName)
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return types.Process{}, err
	}

	message.SetString(action.createdFilePath.Name, filePath)

	return message, nil
}
