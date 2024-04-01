package file

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"os"
	"path/filepath"
)

type contentToFile struct {
	config *configuration.Config
}

func CreateContentToFile(config *configuration.Config) actions.Action {
	return &contentToFile{
		config: config,
	}
}

func (action *contentToFile) Inputs() []actions.Property {
	return []actions.Property{
		{
			Name: "fileName",
			Type: "text",
		},
		{
			Name: "content",
			Type: "text",
		},
	}
}

func (action *contentToFile) Outputs() []actions.Property {
	return []actions.Property{
		{
			Name: "createdFilePath",
			Type: "text",
		},
	}
}

func (action *contentToFile) Execute(message types.Pipeline) (types.Pipeline, error) {
	fileName, err := message.GetString("fileName")
	if err != nil {
		return types.Pipeline{}, err
	}

	content, err := message.GetString("content")
	if err != nil {
		return types.Pipeline{}, err
	}

	filePath := filepath.Join(action.config.Folder.TmpFolder, fileName)
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return types.Pipeline{}, err
	}

	_, _ = message.SetString("createdFilePath", filePath)
	return message, nil
}
