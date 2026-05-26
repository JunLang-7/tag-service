package bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	APP_KEY    = "junlang"
	APP_SECRET = "go-programming-tour-book"
)

type AccessToken struct {
	Token string `json:"token"`
}

type API struct {
	URL    string
	Client *http.Client
}

func NewAPI(url string) *API {
	return &API{url, http.DefaultClient}
}

func (a *API) GetTagList(ctx context.Context, name string) ([]byte, error) {
	token, err := a.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	body, err := a.httpGet(ctx, fmt.Sprintf("%s?token=%s&name=%s", "/api/v1/tags", token, name))
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (a *API) getAccessToken(ctx context.Context) (string, error) {
	body, err := a.httpPost(ctx, "/auth", map[string]string{
		"app_key":    APP_KEY,
		"app_secret": APP_SECRET,
	})
	if err != nil {
		return "", err
	}

	var token AccessToken
	_ = json.Unmarshal(body, &token)
	return token.Token, nil
}

func (a *API) httpGet(ctx context.Context, path string) ([]byte, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s", a.URL, path))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (a *API) httpPost(ctx context.Context, path string, data map[string]string) ([]byte, error) {
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", a.URL, path), bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
