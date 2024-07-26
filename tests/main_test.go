package tests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
)

func TestProcess(t *testing.T) {
	compose, err := tc.NewDockerCompose("../docker/docker-compose.yml")
	require.NoError(t, err, "NewDockerComposeAPI()")

	t.Cleanup(func() {
		require.NoError(t, compose.Down(context.Background(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	require.NoError(t, compose.Up(ctx, tc.Wait(true)), "compose.Up()")

	jsonData := `
	{
    "Parameters": {
        "moje_imie" : "Jan",
        "czesc.text": "Piotr"
    },
    "Tasks": [
        {
            "Sequence": 1,
            "Name": "czesc",
            "Action": "console"
        }
    ]
}
`

	reqBody := bytes.NewBufferString(jsonData)
	//tak wiem...
	time.Sleep(5 * time.Second)
	resp, err := http.Post("http://localhost:5000/execute", "application/json", reqBody)
	if err != nil {
		require.NoError(t, err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	fmt.Println(string(body))

}
