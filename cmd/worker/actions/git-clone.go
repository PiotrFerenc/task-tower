package actions

import (
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/go-git/go-git/v5"
	"os"
)

type gitClone struct {
	config *configuration.Config
}

func CreateGitClone(config *configuration.Config) Action {
	return &gitClone{
		config: config,
	}
}

func (action *gitClone) Inputs() []Property {
	return []Property{
		{
			Name: "url",
			Type: "text",
		},
	}
}

func (action *gitClone) Outputs() []Property {
	return []Property{
		{
			Name: "path",
			Type: "text",
		}}
}
func (action *gitClone) Execute(message types.Pipeline) (types.Pipeline, error) {
	repositoryUrl, err := message.GetString("url")
	if err != nil {
		return types.Pipeline{}, err
	}

	path := message.NewFolder(action.config.Folder.TmpFolder)

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:      repositoryUrl,
		Progress: os.Stdout,
	})

	if err != nil {
		return types.Pipeline{}, err
	}
	_, _ = message.SetString("path", path)
	return message, nil
}
