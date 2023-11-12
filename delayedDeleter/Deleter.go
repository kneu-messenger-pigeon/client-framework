package delayedDeleter

import (
	"context"
	"fmt"
	"github.com/kneu-messenger-pigeon/client-framework/delayedDeleter/contracts"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
	"io"
	"sync"
	"time"
)

const defaultWaitingTimeout = time.Minute

type Deleter struct {
	handler        contracts.DeleteHandlerInterface
	redis          redis.UniversalClient
	out            io.Writer
	queueName      string
	waitingTimeout time.Duration
}

func NewWelcomeAnonymousMessageDelayedDeleter(redis redis.UniversalClient, out io.Writer, name string) *Deleter {
	return &Deleter{
		redis:          redis,
		out:            out,
		queueName:      "deleter_queue_" + name,
		waitingTimeout: defaultWaitingTimeout,
	}
}

func (deleter *Deleter) AddToQueue(messageDeleteTask *contracts.DeleteTask) {
	taskSerialized, _ := proto.Marshal(messageDeleteTask)
	err := deleter.redis.LPush(context.Background(), deleter.queueName, taskSerialized).Err()
	deleter.logError("failed write task to redis queue: ", err)
}

func (deleter *Deleter) SetHandler(handler contracts.DeleteHandlerInterface) {
	deleter.handler = handler
}

func (deleter *Deleter) Execute(ctx context.Context, wg *sync.WaitGroup) {
	var waitTime time.Duration
	var now int64
	var nextTask contracts.DeleteTask
	var taskSerialized []byte
	var nextTasksSerialized []string
	var unmarshalErr error
	var handlerErr error
	var dequeueErr error

	dequeue := func() {
		dequeueErr = deleter.redis.RPop(context.Background(), deleter.queueName).Err()
		deleter.logError("failed dequeue task: ", dequeueErr)
	}

	var readNextTask func()
	readNextTask = func() {
		nextTask.Reset()
		nextTasksSerialized = deleter.redis.LRange(context.Background(), deleter.queueName, -1, -1).Val()
		if len(nextTasksSerialized) > 0 {
			taskSerialized = []byte(nextTasksSerialized[0])
			unmarshalErr = proto.Unmarshal(taskSerialized, &nextTask)
			if unmarshalErr != nil {
				_, _ = fmt.Fprintln(deleter.out, "failed unmarshal task: ", unmarshalErr)
				dequeue()
				readNextTask()
			}
		}
	}

	for {
		now = time.Now().Unix()
		readNextTask()
		for nextTask.GetScheduledAt() != 0 && nextTask.GetScheduledAt() <= now {
			handlerErr = deleter.handler.HandleDeleteTask(&nextTask)
			deleter.logError("handle delete message err: ", handlerErr)
			dequeue()
			readNextTask()
		}

		if nextTask.ScheduledAt > now {
			waitTime = time.Duration(nextTask.GetScheduledAt()-now) * time.Second
		} else {
			waitTime = deleter.waitingTimeout
		}

		select {
		case <-ctx.Done():

			wg.Done()
			return
		case <-time.After(waitTime):
		}
	}
}

func (deleter *Deleter) logError(prefix string, err error) {
	if err != nil {
		_, _ = fmt.Fprintln(deleter.out, prefix, err)
	}
}
