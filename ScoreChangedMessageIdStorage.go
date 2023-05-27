package framework

import (
	"context"
	"fmt"
	"github.com/kneu-messenger-pigeon/client-framework/models"
	"github.com/redis/go-redis/v9"
	"io"
	"strconv"
	"time"
)

type ScoreChangedMessageIdStorage struct {
	out           io.Writer
	redis         redis.UniversalClient
	storageExpire time.Duration
}

func (storage *ScoreChangedMessageIdStorage) Set(studentId uint, lessonId uint, chatId string, messageId string) {
	key := storage.makeKey(studentId, lessonId)
	err := storage.redis.HSet(context.Background(), key, chatId, messageId).Err()

	if err == nil {
		err = storage.redis.Expire(context.Background(), key, storage.storageExpire).Err()
	}

	if err != nil {
		_, _ = fmt.Fprintf(
			storage.out,
			"Failed to save shanged score message id. Error: %s, student %d, lessonId %d, chatId %s, messageId %s",
			err.Error(), studentId, lessonId, chatId, messageId,
		)
	}
}

func (storage *ScoreChangedMessageIdStorage) GetAll(studentId uint, lessonId uint) models.ScoreChangedMessageMap {
	result := storage.redis.HGetAll(context.Background(), storage.makeKey(studentId, lessonId))
	if result.Err() == nil {
		return result.Val()
	}

	return models.ScoreChangedMessageMap{}
}

func (storage *ScoreChangedMessageIdStorage) makeKey(studentId uint, lessonId uint) string {
	return "SM:" + strconv.FormatUint(uint64(studentId), 10) + ":" + strconv.FormatUint(uint64(lessonId), 10)
}
