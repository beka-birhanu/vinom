package service

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	PushToQueue(ctx context.Context, id string, rating int32, latency int32) error
}

type gameSessionManager interface {
	NewSession(ctx context.Context, playerIDs []uuid.UUID)
}

// SortedQueue defines an interface for managing sorted queues.
type sortedQueue interface {
	// Enqueue adds a member to the sorted queue with a given score.
	Enqueue(ctx context.Context, queueKey string, score float64, member string) error

	// DequeTops removes and retrieves up to `amount` members with the lowest scores from the queue.
	DequeTops(ctx context.Context, queueKey string, amount int64) ([]string, error)

	// Count returns the number of members in the sorted queue.
	Count(ctx context.Context, queueKey string) (int64, error)
}
