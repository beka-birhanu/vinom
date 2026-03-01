package service

import (
	"time"

	dmn "github.com/beka-birhanu/vinom/api/domain"
	"github.com/google/uuid"
)

// Tokenizer defines methods for generating and decoding tokens.
type Tokenizer interface {
	// Generate creates a token with the given claims and expiration duration.
	Generate(claims map[string]any, expTime time.Duration) (string, error)

	// Decode validates and parses a token, returning its claims.
	Decode(token string) (map[string]any, error)
}

type Service interface {
	Register(string, string) error
	SignIn(string, string) (*dmn.User, string, error)
}

type UserRepo interface {
	// Save inserts or updates a user in the repository.
	// If the user already exists, it updates the record. Otherwise, it creates a new one.
	Save(user *dmn.User) error

	// ByID retrieves a user by their unique ID.
	// Returns an error if the user is not found or in case of an unexpected error.
	ByID(id uuid.UUID) (*dmn.User, error)

	// ByUsername retrieves a user by their username.
	// Returns an error if the user is not found or in case of an unexpected error.
	ByUsername(username string) (*dmn.User, error)
}
