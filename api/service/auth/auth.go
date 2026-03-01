package service

import (
	"errors"
	"time"

	dmn "github.com/beka-birhanu/vinom/api/domain"
	"github.com/google/uuid"
)

type service struct {
	userRepo  UserRepo
	tokenizer Tokenizer
}

func NewService(ur UserRepo, t Tokenizer) (Service, error) {
	return &service{
		userRepo:  ur,
		tokenizer: t,
	}, nil
}

func (a *service) Register(username, password string) error {
	userConfig := dmn.UserConfig{
		ID:            uuid.New(),
		Username:      username,
		PlainPassword: password,
	}

	_, err := a.userRepo.ByUsername(username)
	if err == nil {
		return errors.New("username already exist")
	}

	user, err := dmn.NewUser(userConfig)
	if err != nil {
		return err
	}

	err = a.userRepo.Save(user)
	if err != nil {
		return err
	}

	return nil
}

func (a *service) SignIn(username, password string) (*dmn.User, string, error) {
	user, err := a.userRepo.ByUsername(username)
	if err != nil {
		return nil, "", errors.New("invalid username or password")
	}

	if !user.VerifyPassword(password) {
		return nil, "", errors.New("invalid username or password")
	}

	token, err := a.tokenizer.Generate(map[string]any{
		"userID":   user.ID,
		"username": user.Username,
	}, 24*time.Hour)

	return user, token, err
}
