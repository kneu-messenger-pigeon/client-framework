package framework

import (
	"bytes"
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redismock/v9"
	"github.com/kneu-messenger-pigeon/client-framework/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScoreChangedMessageIdStorage_Set(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		storage := ScoreChangedMessageIdStorage{
			out:           &bytes.Buffer{},
			redis:         redisClient,
			storageExpire: time.Minute,
		}

		chatId := "test-chat-id-1"
		messageId := "test-message-id"
		expectedKey := "SM:123:99"

		redisMock.ExpectHSet(expectedKey, chatId, messageId).SetVal(1)
		redisMock.ExpectExpire(expectedKey, time.Minute).SetVal(true)

		storage.Set(123, 99, chatId, messageId)

		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("set_empty", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		storage := ScoreChangedMessageIdStorage{
			out:           &bytes.Buffer{},
			redis:         redisClient,
			storageExpire: time.Minute,
		}

		chatId := "test-chat-id-1"
		messageId := ""
		expectedKey := "SM:123:99"

		redisMock.ExpectHDel(expectedKey, chatId).SetVal(1)
		redisMock.ExpectExpire(expectedKey, time.Minute).SetVal(true)

		storage.Set(123, 99, chatId, messageId)

		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("error", func(t *testing.T) {
		expectedError := errors.New("expected error")

		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		out := &bytes.Buffer{}
		storage := ScoreChangedMessageIdStorage{
			out:           out,
			redis:         redisClient,
			storageExpire: time.Minute,
		}

		chatId := "test-chat-id-1"
		messageId := "test-message-id"
		expectedKey := "SM:123:99"

		redisMock.ExpectHSet(expectedKey, chatId, messageId).SetErr(expectedError)

		storage.Set(123, 99, chatId, messageId)

		assert.NoError(t, redisMock.ExpectationsWereMet())
		assert.Contains(t, out.String(), expectedError.Error())
	})
}

func TestScoreChangedMessageIdStorage_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		out := &bytes.Buffer{}
		storage := ScoreChangedMessageIdStorage{
			out:           out,
			redis:         redisClient,
			storageExpire: time.Minute,
		}

		chatId1 := "test-chat-id-1"
		messageId1 := "test-message-id-1"

		chatId2 := "test-chat-id-2"
		messageId2 := "test-message-id-2"
		expectedKey := "SM:123:99"

		expectedResult := models.ScoreChangedMessageMap{
			chatId1: messageId1,
			chatId2: messageId2,
		}

		redisMock.ExpectHGetAll(expectedKey).SetVal(map[string]string{
			chatId1: messageId1,
			chatId2: messageId2,
		})

		actualResult := storage.GetAll(123, 99)

		assert.Equal(t, expectedResult, actualResult)
	})

	t.Run("empty_storage", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		out := &bytes.Buffer{}
		storage := ScoreChangedMessageIdStorage{
			out:           out,
			redis:         redisClient,
			storageExpire: time.Minute,
		}

		expectedKey := "SM:123:99"
		redisMock.ExpectHGetAll(expectedKey).RedisNil()

		actualResult := storage.GetAll(123, 99)
		assert.Equal(t, models.ScoreChangedMessageMap{}, actualResult)
	})

	t.Run("error", func(t *testing.T) {
		expectedError := errors.New("expected error")

		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		out := &bytes.Buffer{}
		storage := ScoreChangedMessageIdStorage{
			out:           out,
			redis:         redisClient,
			storageExpire: time.Minute,
		}

		expectedKey := "SM:123:99"
		redisMock.ExpectHGetAll(expectedKey).SetErr(expectedError)

		actualResult := storage.GetAll(123, 99)
		assert.Equal(t, models.ScoreChangedMessageMap{}, actualResult)
	})
}

func TestScoreChangedMessageIdStorage_SetGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m := miniredis.RunT(t)

		storage := ScoreChangedMessageIdStorage{
			out: &bytes.Buffer{},
			redis: redis.NewClient(&redis.Options{
				Network: "tcp",
				Addr:    m.Addr(),
			}),
			storageExpire: time.Minute,
		}

		studentId := uint(123)
		lessonId := uint(99)

		chatId1 := "test-chat-id-1"
		messageId1 := "test-message-id-1"

		chatId2 := "test-chat-id-2"
		messageId2 := "test-message-id-2"

		expectedResult := models.ScoreChangedMessageMap{
			chatId1: messageId1,
			chatId2: messageId2,
		}

		storage.Set(studentId, lessonId, chatId1, messageId1)
		storage.Set(studentId, lessonId, chatId2, messageId2)

		actualResult := storage.GetAll(studentId, lessonId)

		assert.Equal(t, expectedResult, actualResult)
	})

	t.Run("refresh_ttl_for_for_first_record", func(t *testing.T) {
		m := miniredis.RunT(t)

		storage := ScoreChangedMessageIdStorage{
			out: &bytes.Buffer{},
			redis: redis.NewClient(&redis.Options{
				Network: "tcp",
				Addr:    m.Addr(),
			}),
			storageExpire: time.Minute,
		}

		studentId := uint(123)
		lessonId := uint(99)

		chatId1 := "test-chat-id-1"
		messageId1 := "test-message-id-1"

		chatId2 := "test-chat-id-2"
		messageId2 := "test-message-id-2"

		expectedResult := models.ScoreChangedMessageMap{
			chatId1: messageId1,
			chatId2: messageId2,
		}

		storage.Set(studentId, lessonId, chatId1, messageId1)

		m.FastForward(time.Second * 50)
		storage.Set(studentId, lessonId, chatId2, messageId2)
		m.FastForward(time.Second * 50)

		actualResult := storage.GetAll(studentId, lessonId)

		assert.Equal(t, expectedResult, actualResult)
	})

	t.Run("expire_ttl", func(t *testing.T) {
		m := miniredis.RunT(t)

		storage := ScoreChangedMessageIdStorage{
			out: &bytes.Buffer{},
			redis: redis.NewClient(&redis.Options{
				Network: "tcp",
				Addr:    m.Addr(),
			}),
			storageExpire: time.Minute,
		}

		studentId := uint(123)
		lessonId := uint(99)

		chatId1 := "test-chat-id-1"
		messageId1 := "test-message-id-1"

		chatId2 := "test-chat-id-2"
		messageId2 := "test-message-id-2"

		storage.Set(studentId, lessonId, chatId1, messageId1)
		storage.Set(studentId, lessonId, chatId2, messageId2)
		m.FastForward(time.Minute * 2)

		actualResult := storage.GetAll(studentId, lessonId)

		assert.Equal(t, models.ScoreChangedMessageMap{}, actualResult)
	})
}
