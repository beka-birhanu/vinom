package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	general_i "github.com/beka-birhanu/vinom/common/interfaces/general"
	"github.com/google/uuid"
)

const (
	defaultPrefix           = "matchmaker"
	defaultMaxPlayer        = 2
	defaultRankTolerance    = 0
	defaultLatencyTolerance = 0
	queueRankLatencyKeyFmt  = "%s:queue:rank_%d:latency_%d"
)

var ErrPlayerNotFoundInQueue = errors.New("player not found in queue")

type Options struct {
	Prefix           string
	GameHandler      gameSessionManager
	MaxPlayer        int32
	RankTolerance    int32
	LatencyTolerance int32
}

type service struct {
	sortedQueue sortedQueue
	logger      general_i.Logger
	opts        *Options
}

func NewService(sortedQueue sortedQueue, logger general_i.Logger, opts *Options) (Service, error) {
	if opts == nil {
		opts = &Options{
			MaxPlayer: defaultMaxPlayer,
			Prefix:    defaultPrefix,
		}
	}

	if opts.MaxPlayer <= 0 {
		opts.MaxPlayer = defaultMaxPlayer
	}

	if opts.Prefix == "" {
		opts.Prefix = defaultPrefix
	}

	if opts.RankTolerance < 0 {
		opts.RankTolerance = defaultRankTolerance
	}

	if opts.LatencyTolerance < 0 {
		opts.LatencyTolerance = defaultLatencyTolerance
	}

	return &service{
		opts:        opts,
		sortedQueue: sortedQueue,
		logger:      logger,
	}, nil
}

func (mm *service) PushToQueue(ctx context.Context, id string, rank int32, latency int32) error {
	mm.logger.Info(fmt.Sprintf("Adding player to queue: ID=%s Rank=%d Latency=%d", id, rank, latency))
	score := float64(time.Now().UnixNano())
	err := mm.sortedQueue.Enqueue(ctx, mm.queueKey(rank, latency), score, id)
	if err != nil {
		mm.logger.Error(fmt.Sprintf("Failed to enqueue player: %s", err))
		return err
	}

	mm.logger.Info(fmt.Sprintf("Player enqueued successfully: ID=%s", id))
	go mm.match(context.Background(), rank, latency)
	return nil
}

func (mm *service) match(ctx context.Context, rank int32, latency int32) {
	newCtx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	queueKey := mm.queueKey(rank, latency)
	qLen, err := mm.sortedQueue.Count(newCtx, queueKey)
	if err != nil {
		mm.logger.Error(fmt.Sprintf("Counting queue length: %s", err))
		return
	}

	if qLen >= int64(mm.opts.MaxPlayer) {
		rawPlayers, err := mm.sortedQueue.DequeTops(newCtx, queueKey, int64(mm.opts.MaxPlayer))
		if err != nil {
			mm.logger.Error(fmt.Sprintf("Dequeing tops: %s", err))
			return
		}

		var playersIDs []uuid.UUID
		for _, raw := range rawPlayers {
			if id, err := uuid.Parse(raw); err == nil {
				playersIDs = append(playersIDs, id)
			} else {
				mm.logger.Warning(fmt.Sprintf("Non-UUID value in queue: %s", raw))
			}
		}

		mm.logger.Info(fmt.Sprintf("Match found for players: %v", playersIDs))
		if mm.opts.GameHandler != nil {
			go mm.opts.GameHandler.NewSession(ctx, playersIDs)
		}
	}
}

func (mm *service) SetMatchHandler(gm gameSessionManager) {
	mm.opts.GameHandler = gm
}

func (mm *service) queueKey(rank int32, latency int32) string {
	return fmt.Sprintf(queueRankLatencyKeyFmt, mm.opts.Prefix, scale(rank, mm.opts.RankTolerance), scale(latency, mm.opts.LatencyTolerance))
}

func scale(value, tolerance int32) int32 {
	return value / (tolerance + 1)
}
