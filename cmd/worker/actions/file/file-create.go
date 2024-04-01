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
			Description: "",
			Validation:  "required",
		},
		content: actions.Property{
			Name:        "content",
			Type:        actions.Text,
			Description: "",
			Validation:  "required",
		},
		createdFilePath: actions.Property{
			Name:        "createdFilePath",
			Type:        actions.Text,
			Description: "",
			Validation:  "",
		},
	}
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

func (action *contentToFile) Execute(message types.Pipeline) (types.Pipeline, error) {
	fileName, err := action.fileName.GetStringFrom(&message)
	if err != nil {
		return types.Pipeline{}, err
	}

	content, err := action.content.GetStringFrom(&message)
	if err != nil {
		return types.Pipeline{}, err
	}

	filePath := filepath.Join(action.config.Folder.TmpFolder, fileName)
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return types.Pipeline{}, err
	}

	message.SetString(action.createdFilePath.Name, filePath)

	return message, nil
}
