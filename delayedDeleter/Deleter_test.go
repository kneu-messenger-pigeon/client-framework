package delayedDeleter

import (
	"bytes"
	"context"
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redismock/v9"
	"github.com/kneu-messenger-pigeon/client-framework/delayedDeleter/contracts"
	"github.com/kneu-messenger-pigeon/client-framework/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
	"sync"
	"testing"
	"time"
)

func TestWelcomeAnonymousMessageDelayedDeleter(t *testing.T) {
	matchTask := func(expectedTask *contracts.DeleteTask, actualTask *contracts.DeleteTask) bool {
		if expectedTask == nil && actualTask == nil {
			return true
		}

		return expectedTask.ScheduledAt == actualTask.ScheduledAt &&
			expectedTask.MessageId == actualTask.MessageId &&
			expectedTask.ChatId == actualTask.ChatId
	}

	t.Run("success", func(t *testing.T) {
		task1 := contracts.DeleteTask{
			ScheduledAt: time.Now().Add(time.Second * 2).Unix(),
			MessageId:   1,
			ChatId:      2,
		}

		task2 := contracts.DeleteTask{
			ScheduledAt: time.Now().Add(time.Second * 4).Unix(),
			MessageId:   10,
			ChatId:      20,
		}

		task1FoundCount := 0
		task2FoundCount := 0
		taskNotFoundCount := 0

		handler := mocks.NewDeleteHandlerInterface(t)
		handler.On("HandleDeleteTask", mock.Anything).Times(2).Return(
			func(task *contracts.DeleteTask) error {
				if matchTask(&task1, task) {
					task1FoundCount++
					return nil
				}

				if matchTask(&task2, task) {
					task2FoundCount++
					return nil
				}

				taskNotFoundCount++
				return errors.New("unexpected task")
			},
		)

		out := &bytes.Buffer{}
		redisClient := redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    miniredis.RunT(t).Addr(),
		})

		delayedDeleter := NewWelcomeAnonymousMessageDelayedDeleter(redisClient, out, "test")
		delayedDeleter.waitingTimeout = time.Second
		delayedDeleter.SetHandler(handler)

		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go delayedDeleter.Execute(ctx, wg)

		delayedDeleter.AddToQueue(&task1)
		delayedDeleter.AddToQueue(&task2)
		assert.Empty(t, out.String(), "No error should be printed")

		time.Sleep(time.Second * 5)
		cancel()

		assert.Equal(t, 1, task1FoundCount, "Task 1 should be found once")
		assert.Equal(t, 1, task2FoundCount, "Task 2 should be found once")
		assert.Empty(t, taskNotFoundCount, "No unexpected tasks should be found")
		assert.Empty(t, out.String(), "No error should be printed")
	})

	t.Run("errorUnmarshal", func(t *testing.T) {
		task1 := contracts.DeleteTask{
			ScheduledAt: time.Now().Unix(),
			MessageId:   30,
			ChatId:      40,
		}

		task2 := contracts.DeleteTask{
			ScheduledAt: time.Now().Unix(),
			MessageId:   35,
			ChatId:      45,
		}

		handler := mocks.NewDeleteHandlerInterface(t)
		handler.On("HandleDeleteTask", mock.Anything).Times(2).Return(nil)

		redisClient := redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    miniredis.RunT(t).Addr(),
		})

		out := &bytes.Buffer{}
		delayedDeleter := NewWelcomeAnonymousMessageDelayedDeleter(redisClient, out, "test-errorDequeue")
		delayedDeleter.waitingTimeout = time.Second
		delayedDeleter.SetHandler(handler)

		err := redisClient.LPush(context.Background(), delayedDeleter.queueName, "serializations-errors").Err()
		assert.NoError(t, err)

		delayedDeleter.AddToQueue(&task1)
		delayedDeleter.AddToQueue(&task2)

		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go delayedDeleter.Execute(ctx, wg)

		time.Sleep(time.Second * 2)
		cancel()

		assert.Contains(t, out.String(), "failed unmarshal task: ", "Error should be printed")
	})

	t.Run("errorDequeue", func(t *testing.T) {
		expectedErr := errors.New("dequeue error")
		task1 := contracts.DeleteTask{
			ScheduledAt: time.Now().Unix(),
			MessageId:   50,
			ChatId:      60,
		}
		task1Serialized, _ := proto.Marshal(&task1)

		handler := mocks.NewDeleteHandlerInterface(t)
		handler.On("HandleDeleteTask", mock.Anything).Times(1).Return(nil)

		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		out := &bytes.Buffer{}
		delayedDeleter := NewWelcomeAnonymousMessageDelayedDeleter(redisClient, out, "test-errorDequeue")
		delayedDeleter.waitingTimeout = time.Second
		delayedDeleter.SetHandler(handler)

		redisMock.ExpectLRange(delayedDeleter.queueName, -1, -1).SetVal([]string{
			string(task1Serialized),
		})
		redisMock.ExpectRPop(delayedDeleter.queueName).SetErr(expectedErr)

		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go delayedDeleter.Execute(ctx, wg)

		time.Sleep(time.Second * 2)
		cancel()

		assert.Contains(t, out.String(), "failed dequeue task: ", "Error should be printed")
	})
}
