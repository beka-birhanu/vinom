package service

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/beka-birhanu/vinom-client/dmn"
	"github.com/beka-birhanu/vinom-client/service/i"
)

const (
	loginUri    = "/login"
	registerUri = "/register"
)

type Auth struct {
	httpClient i.HttpRequester
}

// Login implements i.AuthServer.
func (a *Auth) Login(username string, password string) (*dmn.Player, string, error) {
	body := &AuthRequest{
		Username: username,
		Password: password,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, "", err
	}

	response, err := a.httpClient.Post(loginUri, bytes.NewReader(payload))
	if err != nil {
		return nil, "", err
	}

	responseBody, err := io.ReadAll(response)
	if err != nil {
		return nil, "", err
	}

	var loginResponse AuthResponse
	err = json.Unmarshal(responseBody, &loginResponse)
	if err != nil {
		return nil, "", err // Return error if unmarshalling fails
	}

	// Return the player, token, and nil error
	return &dmn.Player{
		ID:       loginResponse.ID,
		Rating:   loginResponse.Rating,
		Username: loginResponse.Username,
	}, loginResponse.Token, nil
}

// Register implements i.AuthServer.
func (a *Auth) Register(string, string) error {
	panic("unimplemented")
}

func NewAuth(hr i.HttpRequester) (i.AuthServer, error) {
	return &Auth{
		httpClient: hr,
	}, nil
}
