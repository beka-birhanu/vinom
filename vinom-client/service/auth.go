package service

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/beka-birhanu/vinom-client/dmn"
	"github.com/beka-birhanu/vinom-client/service/i"
)

type Auth struct {
	httpClient  i.HttpRequester
	loginUri    string
	registerUri string
}

func NewAuth(hr i.HttpRequester, loginUri, registerUri string) (i.AuthServer, error) {
	return &Auth{
		httpClient:  hr,
		loginUri:    loginUri,
		registerUri: registerUri,
	}, nil
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

	response, err := a.httpClient.Post(a.loginUri, bytes.NewReader(payload), "")
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
func (a *Auth) Register(username string, password string) error {
	body := &AuthRequest{
		Username: username,
		Password: password,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	_, err = a.httpClient.Post(a.registerUri, bytes.NewReader(payload), "")
	if err != nil {
		return err
	}

	return nil
}
