package client

import (
	"bytes"
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

var (
	userAgent = "todo-app-client"
)

// Get is a generic http function for unmarshalling a request to json
func Get[M any](client Client, url string) (*M, error) {
	return getOrDelete[M](client, "GET", url)
}

// Delete is a generic http function for deleting a resource and unmarshalling
// the response to json.
func Delete[M any](client Client, url string) (*M, error) {
	return getOrDelete[M](client, "DELETE", url)
}

// Post is a generic http function for creating a resource and unmarshalling
// the response to json.
func Post[M any](client Client, url string, m *M) (*M, error) {
	return postOrPut(client, "POST", url, m)
}

// Put is a generic http function for updating a resource and unmarshalling
// the response to json.
func Put[M any](client Client, url string, m *M) (*M, error) {
	return postOrPut(client, "PUT", url, m)
}

func getOrDelete[M any](client Client, verb, url string) (*M, error) {
	url = strings.TrimPrefix(url, "/")
	req, err := http.NewRequest(verb, client.BaseUrl+url, nil)
	if err != nil {
		return nil, err
	}

	// TODO: would setting agent and bearer go better in a round tripper?
	req.Header.Set("User-Agent", userAgent)

	token, err := client.Tokens.GetIdToken()
	if err != nil {
		return nil, err
	}
	authHeader := fmt.Sprintf("BEARER %s", token)
	req.Header.Set("Authorization", authHeader)

	resp, err := client.Http.Do(req)

	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
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

// Get is a generic http function for unmarshalling a request to json
func postOrPut[M any](client Client, verb, url string, m *M) (*M, error) {
	url = strings.TrimPrefix(url, "/")

	postData, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(verb, client.BaseUrl+url, bytes.NewBuffer(postData))
	if err != nil {
		return nil, err
	}

	// TODO: would setting agent and bearer go better in a round tripper?
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

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

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.New("Non 200 response: " + resp.Status + " " + string(data))
	}

	var model M
	if err := json.Unmarshal(data, &model); err != nil {
		return nil, err
	}

	return &model, nil
}
