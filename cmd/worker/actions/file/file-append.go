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
			Description: "",
			Validation:  "required",
		},
		content: actions.Property{
			Name:        "content",
			Type:        actions.Text,
			Description: "",
			Validation:  "required",
		},
		appendFilePath: actions.Property{
			Name:        "appendFilePath",
			Type:        actions.Text,
			Description: "",
		},
	}
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

func (action *appendContentToFile) Execute(message types.Pipeline) (types.Pipeline, error) {
	fileName, err := action.fileName.GetStringFrom(&message)
	if err != nil {
		return types.Pipeline{}, err
	}
	content, err := action.content.GetStringFrom(&message)
	if err != nil {
		return types.Pipeline{}, err
	}
	filePath := filepath.Join(action.config.Folder.TmpFolder, fileName)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return types.Pipeline{}, err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return types.Pipeline{}, err
	}
	message.SetString(action.appendFilePath.Name, filePath)
	return message, nil
}
