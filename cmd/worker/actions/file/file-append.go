package file

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"os"
	"path/filepath"
)

type appendContentToFile struct {
	config         *configuration.Config
	fileName       actions.Property
	content        actions.Property
	appendFilePath actions.Property
}

func CreateAppendContentToFile(config *configuration.Config) actions.Action {
	return &appendContentToFile{
		config: config,
		fileName: actions.Property{
			Name:        "fileName",
			Type:        actions.Text,
			Description: "The name of the file to which content is being appended",
			DisplayName: "File Name",
			Validation:  "required",
		},
		content: actions.Property{
			Name:        "content",
			Type:        actions.Text,
			Description: "The content to be appended to the file",
			DisplayName: "Content",
			Validation:  "required",
		},
		appendFilePath: actions.Property{
			Name:        "appendFilePath",
			Type:        actions.Text,
			Description: "The path of the file to which the content will be appended",
			DisplayName: "Append File Path",
			Validation:  "",
		},
	}
}

func (action *appendContentToFile) GetCategoryName() string {
	return "file"
}
func (action *appendContentToFile) Inputs() []actions.Property {
	return []actions.Property{
		action.fileName, action.content,
	}
}

func (action *appendContentToFile) Outputs() []actions.Property {
	return []actions.Property{
		action.appendFilePath,
	}
}

func (action *appendContentToFile) Execute(message types.Process) (types.Process, error) {
	fileName, err := action.fileName.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}
	content, err := action.content.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}
	filePath := filepath.Join(action.config.Folder.TmpFolder, fileName)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return types.Process{}, err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return types.Process{}, err
	}
	message.SetString(action.appendFilePath.Name, filePath)
	return message, nil
}
