package git

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/go-git/go-git/v5"
	"os"
)

type gitClone struct {
	config *configuration.Config
	url    actions.Property
	path   actions.Property
}

func CreateGitClone(config *configuration.Config) actions.Action {
	return &gitClone{
		config: config,
		url: actions.Property{
			Name:        "url",
			Type:        actions.Text,
			Description: "URL of the git repository to clone.",
			DisplayName: "Git Repository URL",
			Validation:  "required",
		},
		path: actions.Property{
			Name:        "path",
			Type:        actions.Text,
			Description: "Local path where the repository will be cloned to.",
			DisplayName: "Destination Path",
			Validation:  "required",
		},
	}
}

func (action *gitClone) GetCategoryName() string {
	return "git"
}
func (action *gitClone) Inputs() []actions.Property {
	return []actions.Property{
		action.url,
	}
}

func (action *gitClone) Outputs() []actions.Property {
	return []actions.Property{
		action.path,
	}
}
func (action *gitClone) Execute(message types.Process) (types.Process, error) {
	repositoryUrl, err := action.url.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}

	path := message.NewFolder(action.config.Folder.TmpFolder)

	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:      repositoryUrl,
		Progress: os.Stdout,
	})

	if err != nil {
		return types.Process{}, err
	}
	message.SetString(action.path.Name, path)
	return message, nil
}
