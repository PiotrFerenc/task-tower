package git

import (
	"github.com/PiotrFerenc/mash2/cmd/worker/actions"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"time"
)

type gitCommit struct {
	config  *configuration.Config
	path    actions.Property
	message actions.Property
	id      actions.Property
}

func CreateGitCommit(config *configuration.Config) actions.Action {
	return &gitCommit{
		config: config,
		path: actions.Property{
			Name:        "path",
			Type:        actions.Text,
			Description: "Path to the Git repository",
			DisplayName: "Repository Path",
			Validation:  "required",
		},
		message: actions.Property{
			Name:        "message",
			Type:        actions.Text,
			Description: "Message for the Git commit",
			DisplayName: "Commit Message",
			Validation:  "required",
		},
		id: actions.Property{
			Name:        "id",
			Type:        actions.Text,
			Description: "ID of the Git commit",
			DisplayName: "Commit ID",
			Validation:  "required",
		},
	}
}

func (action *gitCommit) GetCategoryName() string {
	return "git"
}
func (action *gitCommit) Inputs() []actions.Property {
	return []actions.Property{
		action.path,
		action.message,
	}
}

func (action *gitCommit) Outputs() []actions.Property {
	return []actions.Property{
		action.id,
	}
}

// Execute receives a message of type types.Process and performs a Git commit operation.
// It retrieves the repository path and commit message from the message, opens the repository, adds all files, and commits the changes with the provided commit options.
// It returns the modified message or an error if any operation fails.
func (action *gitCommit) Execute(message types.Process) (types.Process, error) {
	repoPath, err := action.path.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}
	commitMessage, err := action.message.GetStringFrom(&message)
	if err != nil {
		return types.Process{}, err
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return types.Process{}, err
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return types.Process{}, err
	}

	_, err = worktree.Add(".")
	if err != nil {
		return types.Process{}, err
	}

	_, _ = worktree.Commit(commitMessage, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Your Name",
			Email: "you@example.com",
			When:  time.Now(),
		},
	})

	if err != nil {
		return types.Process{}, err
	}

	//	object, err := repo.CommitObject(plumbing.Hash{commit})

	if err != nil {
		return types.Process{}, err
	}

	return message, nil
}
