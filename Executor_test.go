package framework

import (
	"bytes"
	"context"
	"github.com/go-redis/redismock/v9"
	frameworkMocks "github.com/kneu-messenger-pigeon/client-framework-mocks"
	"github.com/stretchr/testify/mock"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestEventLoopExecute(t *testing.T) {
	t.Run("Executor Execute", func(t *testing.T) {
		matchContext := mock.MatchedBy(func(ctx context.Context) bool { return true })
		matchWaitGroup := mock.MatchedBy(func(wg *sync.WaitGroup) bool { wg.Done(); return true })

		processor := frameworkMocks.NewMockKafkaConsumerProcessorInterface(t)
		processor.On("Execute", matchContext, matchWaitGroup).Return().Times(2)

		clientController := frameworkMocks.NewMockClientControllerInterface(t)
		clientController.On("Execute", matchContext, matchWaitGroup).Return().Times(1)

		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		redisMock.ExpectPing()
		redisMock.ExpectSave()

		executor := Executor{
			out: &bytes.Buffer{},
			serviceContainer: &ServiceContainer{
				UserRepository: &UserRepository{
					out:   &bytes.Buffer{},
					redis: redisClient,
				},

				UserAuthorizedEventProcessor: processor,
				ScoreChangedEventProcessor:   processor,
				ClientController:             clientController,
			},
		}

		go func() {
			time.Sleep(time.Nanosecond * 200)
			_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()
		executor.Execute()

		processor.AssertExpectations(t)
		clientController.AssertExpectations(t)
	})
}
