package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/csams/doit/pkg/auth"
)

func Get[M any](client *http.Client, url string, tokenProvider *auth.TokenProvider) (*M, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// TODO: convert to configuration
	req.Header.Set("User-Agent", "todo-app-client")

	token, err := tokenProvider.GetIdToken()
	if err != nil {
		return nil, err
	}
	authHeader := fmt.Sprintf("BEARER %s", token)
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var model M
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return &model, nil
}
