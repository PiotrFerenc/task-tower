package file

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"os"
	"path/filepath"
)

type deleteFile struct {
	config   *configuration.Config
	fileName actions.Property
}

func CreateDeleteFileAction(config *configuration.Config) actions.Action {
	return &deleteFile{
		config: config,
		fileName: actions.Property{
			Name:        "fileName",
			Type:        actions.Text,
			Description: "Name of the file to be deleted",
			DisplayName: "File Name",
			Validation:  "required",
		},
	}
}

func (action *deleteFile) GetCategoryName() string {
	return "file"
}
func (action *deleteFile) Inputs() []actions.Property {
	return []actions.Property{
		action.fileName,
	}
}

func (action *deleteFile) Outputs() []actions.Property {
	return []actions.Property{}
}

func (action *deleteFile) Execute(message types.Pipeline) (types.Pipeline, error) {
	fileName, err := action.fileName.GetStringFrom(&message)
	if err != nil {
		return types.Pipeline{}, err
	}
	filePath := filepath.Join(action.config.Folder.TmpFolder, fileName)
	err = os.Remove(filePath)
	return message, err
}
