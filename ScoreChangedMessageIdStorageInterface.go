package framework

import "github.com/kneu-messenger-pigeon/client-framework/models"

type ScoreChangedMessageIdStorageInterface interface {
	Set(studentId uint, lessonId uint, chatId string, messageId string)
	GetAll(studentId uint, lessonId uint) models.ScoreChangedMessageMap
}
