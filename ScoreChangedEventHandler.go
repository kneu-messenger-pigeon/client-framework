package framework

import (
	"fmt"
	"github.com/kneu-messenger-pigeon/events"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	"github.com/kneu-messenger-pigeon/score-client"
	"io"
)

type ScoreChangedEventHandler struct {
	out                          io.Writer
	serviceContainer             *ServiceContainer
	debugLogger                  *DebugLogger
	repository                   UserRepositoryInterface
	scoreClient                  score.ClientInterface
	scoreChangedEventComposer    ScoreChangeEventComposerInterface
	scoreChangedStateStorage     ScoreChangedStateStorageInterface
	scoreChangedMessageIdStorage ScoreChangedMessageIdStorageInterface
	multiMutex                   MultiMutex
}

type ScoreChangedEventPayload struct {
	scoreApi.DisciplineScore
	Previous scoreApi.Score
}

func (handler *ScoreChangedEventHandler) GetExpectedMessageKey() string {
	return events.ScoreChangedEventName
}

func (handler *ScoreChangedEventHandler) GetExpectedEventType() any {
	return &events.ScoreChangedEvent{}
}

func (handler *ScoreChangedEventHandler) Commit() error {
	return nil
}

func (handler *ScoreChangedEventHandler) Handle(s any) error {
	event := s.(*events.ScoreChangedEvent)

	if handler.serviceContainer == nil || handler.serviceContainer.ClientController == nil {
		return nil
	}

	chatIds := handler.repository.GetClientUserIds(event.StudentId)

	handler.debugLogger.Log(
		"ScoreChangedEventHandler: receive change lessonId %d, studentId %d, chatIds: %v",
		event.LessonId, event.StudentId, chatIds,
	)
	if len(chatIds) == 0 {
		return nil
	}

	disciplineScore, err := handler.scoreClient.GetStudentScore(
		uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId),
	)
	if err != nil {
		return err
	}

	go handler.callControllerAction(event, &chatIds, &disciplineScore)
	return nil
}

func (handler *ScoreChangedEventHandler) callControllerAction(
	event *events.ScoreChangedEvent, chatIds *[]string,
	disciplineScore *scoreApi.DisciplineScore,
) {
	mutex := handler.multiMutex.Get(event.Id)
	mutex.Lock()
	defer mutex.Unlock()

	previousScore := handler.scoreChangedEventComposer.Compose(event, &disciplineScore.Score)

	newState := CalculateState(&disciplineScore.Score, previousScore)
	if handler.scoreChangedStateStorage.Get(event.StudentId, event.LessonId) == newState {
		handler.debugLogger.Log(
			"ScoreChangedEventHandler: change (lessonId:%d, studentId %d) skip as not changed",
			event.LessonId, event.StudentId,
		)
		return
	}

	previousMessageIds := handler.scoreChangedMessageIdStorage.GetAll(event.StudentId, event.LessonId)
	for _, chatId := range *chatIds {
		err, newMessageId := handler.serviceContainer.ClientController.ScoreChangedAction(
			chatId, previousMessageIds[chatId], disciplineScore, previousScore,
		)

		handler.debugLogger.Log(
			"ScoreChangedEventHandler: change (lessonId:%d, studentId %d) send update: message id %d",
			event.LessonId, event.StudentId, newMessageId,
		)

		if err != nil {
			_, _ = fmt.Fprintf(handler.out, "ScoreChangedAction return error: %v\n", err)
		}

		if newMessageId != previousMessageIds[chatId] && (newMessageId != "" || err == nil) {
			handler.scoreChangedMessageIdStorage.Set(event.StudentId, event.LessonId, chatId, newMessageId)
		}
	}

	handler.scoreChangedStateStorage.Set(event.StudentId, event.LessonId, newState)
}
