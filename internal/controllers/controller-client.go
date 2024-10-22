package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/PiotrFerenc/mash2/internal/configuration"
	"github.com/google/uuid"
	"io"
	"net/http"
)

type ControllerClient interface {
	Execute(pipeline types.Pipeline) (uuid.UUID, error)
}
type client struct {
	config configuration.ControllerConfig
}

// CreateControllerClient creates a ControllerClient using the provided ControllerConfig.
func CreateControllerClient(config configuration.ControllerConfig) ControllerClient {
	return &client{
		config: config,
	}

}

// Execute executes the given pipeline by sending a POST request to the specified URL with the pipeline data.
// It returns the process ID if the execution is successful, otherwise an error is returned.
func (c *client) Execute(pipeline types.Pipeline) (uuid.UUID, error) {
	data, err := json.Marshal(pipeline)
	if err != nil {
		return uuid.Nil, err
	}
	url := fmt.Sprintf("%s:%s/execute", c.config.Host, c.config.Post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return uuid.Nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return uuid.Nil, err
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return uuid.Nil, err
	}

	var processResponse ProcessResponse
	jsonErr := json.Unmarshal(body, &processResponse)
	if jsonErr != nil {
		return uuid.Nil, err
	}
	return processResponse.ProcessId, nil
}
