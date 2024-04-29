package editor

import (
	"encoding/json"
	"fmt"
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/PiotrFerenc/mash2/internal/executor"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateEditorHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		initCommand, err := initCommand()
		if err != nil {
			return err
		}
		initActions, err := initActions()
		if err != nil {
			return err
		}
		data := map[string]interface{}{
			"Title":       "Edytor",
			"InitCommand": string(initCommand),
			"Actions":     initActions,
		}
		return c.Render(http.StatusOK, "edytor.html", data)
	}
}
func ExecutePipelineHandler() func(c echo.Context) error {
	return func(c echo.Context) error {

	}
}
func initActions() (map[string]string, error) {
	actions := executor.CreateActionMap(&configuration.Config{})
	result := make(map[string]string, len(actions))

	for name, _ := range actions {
		jsonString, err := json.MarshalIndent(types.Task{
			Sequence: 0,
			Action:   name,
			Name:     fmt.Sprintf("my-%s", name),
		}, "", "	")
		if err != nil {
			return nil, err
		}
		result[name] = string(jsonString)

	}
	return result, nil
}

func initCommand() ([]byte, error) {
	initCommand := types.Pipeline{
		Parameters: map[string]interface{}{
			"my-console.text": "hallo word",
		},
		Tasks: []types.Task{
			types.Task{
				Sequence: 1,
				Action:   "console",
				Name:     "my-console",
			},
		},
	}
	initCommandBytes, err := json.MarshalIndent(initCommand, "", "	")
	if err != nil {
		return nil, err
	}
	return initCommandBytes, nil
}
