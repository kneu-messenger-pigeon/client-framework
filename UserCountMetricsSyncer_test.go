package framework

import (
	"bytes"
	"context"
	"github.com/kneu-messenger-pigeon/client-framework/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestUserCountMetricsSyncer_Execute(t *testing.T) {
	expectCount := uint64(221)

	userRepository := mocks.NewUserRepositoryInterface(t)
	userRepository.On("GetUserCount", mock.Anything).Return(expectCount, nil).Times(2)

	userCountMetricsSyncer := UserCountMetricsSyncer{
		out:            &bytes.Buffer{},
		userRepository: userRepository,
	}

	userSyncInterval = time.Second

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go userCountMetricsSyncer.Execute(ctx, wg)

	runtime.Gosched()
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, expectCount, userCount.Get())

	time.Sleep(userSyncInterval)
	cancel()

	wg.Wait()

	assert.Equal(t, expectCount, userCount.Get())
	userRepository.AssertNumberOfCalls(t, "GetUserCount", 2)
}
