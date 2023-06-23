package framework

import "sync"

// MultiMutexCount calculated:
/**
 * Scores-changes-feed partitions count = 6;
 * Telegram rate limit 30 messages / seconds;
 * 30 / 6 = 5; make it double bigger because of hash collusion and not evenly load balancing
 */
const MultiMutexCount = uint(10)

// MultiMutex is a fixed length structure of sync.Mutex
type MultiMutex [MultiMutexCount]sync.Mutex

// Get retrieves a sync.Mutex from an interface
func (m *MultiMutex) Get(key uint) *sync.Mutex {
	return &m[key%MultiMutexCount]
}
