package framework

import (
	"context"
	"fmt"
	"io"
	"os/signal"
	"sync"
	"syscall"
)

type Executor struct {
	out              io.Writer
	serviceContainer *ServiceContainer
}

func (executor *Executor) Execute() {
	_, err := executor.serviceContainer.UserRepository.redis.Ping(context.Background()).Result()
	if err != nil {
		_, _ = fmt.Fprintf(executor.out, "Failed to connect to redisClient: %s\n", err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	wg := &sync.WaitGroup{}

	wg.Add(4)
	go executor.serviceContainer.ClientController.Execute(ctx, wg)
	go executor.serviceContainer.UserAuthorizedEventProcessor.Execute(ctx, wg)
	go executor.serviceContainer.UserCountMetricsSyncer.Execute(ctx, wg)
	go executor.serviceContainer.WelcomeAnonymousDelayedDeleter.Execute(ctx, wg)

	wg.Add(len(executor.serviceContainer.ScoreChangedEventProcessorPool))
	for _, scoreChangedEventProcessor := range executor.serviceContainer.ScoreChangedEventProcessorPool {
		go scoreChangedEventProcessor.Execute(ctx, wg)
	}

	wg.Wait()

	executor.serviceContainer.UserRepository.redis.Save(context.Background())
}
