package editor

import (
	"encoding/json"
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/labstack/echo/v4"
	"net/http"
)

func CreateEditorHandler() func(c echo.Context) error {
	return func(c echo.Context) error {
		initCommand, err := initCommand()
		if err != nil {
			return err
		}

		data := map[string]interface{}{
			"Title":       "Edytor",
			"InitCommand": string(initCommand),
		}
		return c.Render(http.StatusOK, "edytor.html", data)
	}
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
