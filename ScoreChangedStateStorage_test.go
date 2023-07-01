package framework

import (
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScoreChangedStateStorage_Get(t *testing.T) {
	studentId := uint(519)
	lessonId := uint(678)

	expectedState := "12s923z"
	expectedKey := "S519:678"

	t.Run("simple", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)
		redisMock.ExpectGetEx(expectedKey, time.Minute).SetVal(expectedState)

		storage := ScoreChangedStateStorage{
			redis:         redisClient,
			storageExpire: time.Minute,
		}

		actualState := storage.Get(studentId, lessonId)
		assert.NoError(t, redisMock.ExpectationsWereMet())
		assert.Equal(t, expectedState, actualState)
	})
}

func TestScoreChangedStateStorage_Set(t *testing.T) {
	studentId := uint(519)
	lessonId := uint(678)

	expectedState := "12s923z"
	expectedKey := "S519:678"

	t.Run("simple", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)
		redisMock.ExpectSet(expectedKey, expectedState, time.Minute).SetVal("OK")

		storage := ScoreChangedStateStorage{
			redis:         redisClient,
			storageExpire: time.Minute,
		}

		storage.Set(studentId, lessonId, expectedState)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})
}

func TestScoreChangedStateStorage_GetSet(t *testing.T) {
	studentId := uint(519)
	lessonId := uint(678)
	expectedState := "12s923z"

	t.Run("simple", func(t *testing.T) {
		storage := ScoreChangedStateStorage{
			redis: redis.NewClient(&redis.Options{
				Network: "tcp",
				Addr:    miniredis.RunT(t).Addr(),
			}),
			storageExpire: time.Minute,
		}

		storage.Set(studentId, lessonId, expectedState)

		actualState := storage.Get(studentId, lessonId)
		assert.Equal(t, expectedState, actualState)
	})
}
