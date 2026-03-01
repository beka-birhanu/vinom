package service

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// GameService defines the interface for a maze game service.
type GameService interface {
	// Start begins the game and listens for player actions or a timeout.
	Start(ctx context.Context, timeout time.Duration)

	// Stop ends the game, closes channels, and broadcasts the final state.
	Stop()

	// StateChan returns the state change channel.
	StateChan() <-chan []byte

	// ActionChan returns the action channel.
	ActionChan() chan<- []byte

	// EndChan returns the end channel for the game.
	EndChan() <-chan []byte
}

// SessionService manages game sessions and provides session-related information.
type SessionService interface {
	NewSession(ctx context.Context, playerIDs []uuid.UUID)
	StopAll()
	SessionInfo(ctx context.Context, id uuid.UUID) (publicKey []byte, serverAddr string, err error)
}
