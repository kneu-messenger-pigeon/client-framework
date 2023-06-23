package framework

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestFixed_Get(t *testing.T) {
	key := uint(1)
	multiMutex := MultiMutex{}
	m := multiMutex.Get(key)
	m.Lock()
	defer m.Unlock()
	assert.Equal(t, &multiMutex[key], m)
	assert.Equal(t, &multiMutex[key], multiMutex.Get(key))
}

func TestFixed_Get_One(t *testing.T) {
	key1 := uint(1)
	key2 := key1 + MultiMutexCount
	multiMutex := MultiMutex{}
	m := multiMutex.Get(key1)
	m.Lock()
	defer m.Unlock()
	assert.Equal(t, &multiMutex[key1], m)
	assert.Equal(t, &multiMutex[key1], multiMutex.Get(key2))
}

func TestFixed_Get_Race(t *testing.T) {
	key1 := uint(1)
	key2 := key1 + MultiMutexCount

	wg := sync.WaitGroup{}
	wg.Add(2)
	multiMutex := MultiMutex{}

	counter := int32(0)

	go func() {
		m := multiMutex.Get(key1)
		m.Lock()
		atomic.AddInt32(&counter, 1)
		time.Sleep(time.Millisecond * 100)
		atomic.AddInt32(&counter, -1)
		defer m.Unlock()
		wg.Done()
	}()

	go func() {
		m := multiMutex.Get(key2)
		m.Lock()
		assert.Equal(t, int32(0), counter)
		defer m.Unlock()
		wg.Done()
	}()

	wg.Wait()
}
