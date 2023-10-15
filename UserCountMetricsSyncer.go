package framework

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"
)

var userSyncInterval = time.Hour * 8

type UserCountMetricsSyncerInterface interface {
	Execute(ctx context.Context, wg *sync.WaitGroup)
}

type UserCountMetricsSyncer struct {
	out            io.Writer
	userRepository UserRepositoryInterface
}

func (s *UserCountMetricsSyncer) Execute(ctx context.Context, wg *sync.WaitGroup) {
	ticker := time.NewTicker(userSyncInterval)

syncerLoop:
	for {
		repositoryUserCount, err := s.userRepository.GetUserCount(ctx)
		_, _ = fmt.Fprintf(
			s.out,
			"UserCountMetricsSyncer: %d (err: %v) \n", repositoryUserCount, err,
		)

		if err == nil {
			userCount.Set(repositoryUserCount)
		}

		select {
		case <-ctx.Done():
			break syncerLoop
		case <-ticker.C:
		}
	}

	wg.Done()
}
