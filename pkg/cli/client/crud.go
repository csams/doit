package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/csams/doit/pkg/auth"
)

type Client struct {
	Http    *http.Client
	Tokens  *auth.TokenProvider
	BaseUrl string
}

func NewClient(h *http.Client, t *auth.TokenProvider, b string) Client {
	return Client{
		Http:    h,
		Tokens:  t,
		BaseUrl: b,
	}
}

// Get is a generic http function for unmarshalling a request to json
func Get[M any](client Client, url string) (*M, error) {
	url = strings.TrimPrefix(url, "/")
	req, err := http.NewRequest("GET", client.BaseUrl+url, nil)
	if err != nil {
		return nil, err
	}

	// TODO: would setting agent and bearer go better in a round tripper?
	req.Header.Set("User-Agent", "todo-app-client")

	token, err := client.Tokens.GetIdToken()
	if err != nil {
		return nil, err
	}
	authHeader := fmt.Sprintf("BEARER %s", token)
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Http.Do(req)
	if resp.Body != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New("Non 200 response: " + resp.Status)
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
