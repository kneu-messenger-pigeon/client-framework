package framework

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type ScoreChangedStateStorage struct {
	redis         redis.UniversalClient
	storageExpire time.Duration
}

func (storage *ScoreChangedStateStorage) Set(studentId uint, lessonId uint, state string) {
	storage.redis.Set(context.Background(), storage.makeKey(studentId, lessonId), state, storage.storageExpire)
}

func (storage *ScoreChangedStateStorage) Get(studentId uint, lessonId uint) string {
	return storage.redis.GetEx(context.Background(), storage.makeKey(studentId, lessonId), storage.storageExpire).Val()
}

func (storage *ScoreChangedStateStorage) makeKey(studentId uint, lessonId uint) string {
	return "S" + strconv.FormatUint(uint64(studentId), 10) + ":" + strconv.FormatUint(uint64(lessonId), 10)
}
