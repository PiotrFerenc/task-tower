package git

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
)

type gitCreateBranch struct {
	config     *configuration.Config
	repoPath   actions.Property
	branchName actions.Property
}

func CreateGitCreateBranch(config *configuration.Config) actions.Action {
	return &gitCreateBranch{
		config: config,
		repoPath: actions.Property{
			Name:        "repoPath",
			Type:        actions.Text,
			Description: "The path to the local Git repository where the new branch will be created.",
			DisplayName: "Repository Path",
			Validation:  "required",
		},
		branchName: actions.Property{
			Name:        "branchName",
			Type:        actions.Text,
			Description: "The name of the new branch to create in the Git repository.",
			DisplayName: "Branch Name",
			Validation:  "required",
		},
	}
}

func (action *gitCreateBranch) GetCategoryName() string {
	return "git"
}
func (action *gitCreateBranch) Inputs() []actions.Property {
	return []actions.Property{
		action.repoPath,
		action.branchName,
	}
}

// Outputs method returns an empty list of actions.Property.
// It is a method of gitCreateBranch struct, which is a implementation of the Action interface.
// The function is used to specify the output properties after executing the action.
// It returns an empty slice of actions.Property.
func (action *gitCreateBranch) Outputs() []actions.Property {
	return []actions.Property{}
}

func (action *gitCreateBranch) Execute(message types.Process) (types.Process, error) {
	repoPath, _ := action.repoPath.GetStringFrom(&message)
	branchName, _ := action.branchName.GetStringFrom(&message)
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return types.Process{}, err
	}

	err = r.CreateBranch(&config.Branch{
		Name: branchName,
	})
	return message, err
}
