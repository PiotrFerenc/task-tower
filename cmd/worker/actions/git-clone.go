package actions

import (
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/go-git/go-git/v5"
	"os"
)

type gitClone struct {
}

func CreateGitClone() Action {
	return &gitClone{}
}

func (action *gitClone) Inputs() []Property {
	output := make([]Property, 1)
	output[0] = Property{
		Name: "url",
		Type: "text",
	}
	return output
}
func (action *gitClone) Outputs() []Property {
	output := make([]Property, 1)
	output[0] = Property{
		Name: "path",
		Type: "text",
	}
	return output
}
func (action *gitClone) Execute(message types.Message) (types.Message, error) {
	repositoryUrl, err := message.GetString("url")
	if err != nil {
		return types.Message{}, err
	}

	path := "/tmp/"
	_, err = git.PlainClone(path, false, &git.CloneOptions{
		URL:      repositoryUrl,
		Progress: os.Stdout,
	})

	if err != nil {
		return types.Message{}, err
	}
	_, _ = message.SetString("path", path)
	return message, nil
}
