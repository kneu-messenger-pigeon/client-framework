package main

import (
	"context"
	"sync"
)

type ExecutableInterface interface {
	Execute(ctx context.Context, wg *sync.WaitGroup)
}
