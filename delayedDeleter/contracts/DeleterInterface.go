package contracts

import (
	"context"
	"sync"
)

type DeleterInterface interface {
	AddToQueue(messageDeleteTask *DeleteTask)
	Execute(ctx context.Context, wg *sync.WaitGroup)
	SetHandler(handler DeleteHandlerInterface)
}
