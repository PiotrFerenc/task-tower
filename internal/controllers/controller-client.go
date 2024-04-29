package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/PiotrFerenc/mash2/api/types"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
)

type ControllerClient interface {
	Execute(url string, pipeline types.Pipeline) (uuid.UUID, error)
}
type client struct {
}

func CreateControllerClient() ControllerClient {
	return &client{}
}
func (c *client) Execute(url string, pipeline types.Pipeline) (uuid.UUID, error) {

	data := []byte(`{"key":"value"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalf("Error occurred while creating request: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error occurred while sending request: %s", err.Error())
	}

	defer resp.Body.Close()

	body, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatalf("Error occurred while reading response: %s", readErr.Error())
	}

	var processResponse ProcessResponse
	jsonErr := json.Unmarshal(body, &processResponse)
	if jsonErr != nil {
		log.Fatalf("Error occurred while unmarshalling response: %s", jsonErr.Error())
	}
	return processResponse.ProcessId, nil
}
