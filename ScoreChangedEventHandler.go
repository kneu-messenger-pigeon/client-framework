package framework

import (
	"fmt"
	"github.com/kneu-messenger-pigeon/events"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	"github.com/kneu-messenger-pigeon/score-client"
	"io"
	"time"
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
	waitingForAnotherScoreTime   time.Duration
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

	handler.scoreChangedEventComposer.SavePreviousScore(event)

	go handler.callControllerAction(event, &chatIds)
	return nil
}

func (handler *ScoreChangedEventHandler) callControllerAction(
	event *events.ScoreChangedEvent, chatIds *[]string,
) {
	if !handler.scoreChangedEventComposer.BothPreviousScoresSaved(event) {
		time.Sleep(handler.waitingForAnotherScoreTime)
	}

	// each score could be identified by lessonId and studentId
	mutex := handler.multiMutex.Get((event.LessonId & event.StudentId) * (event.LessonId | event.StudentId))
	mutex.Lock()
	defer mutex.Unlock()
	handler.debugLogger.Log("Get lock to process event: %v", event)

	disciplineScore, err := handler.scoreClient.GetStudentScore(
		uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId),
	)
	if err != nil {
		_, _ = fmt.Fprintf(handler.out, "GetStudentScore return error: %v\n", err)
		return
	}

	previousScore := handler.scoreChangedEventComposer.Compose(event, &disciplineScore.Score)

	handler.debugLogger.Log(
		"ScoreChangedEventHandler: change (lessonId:%d, studentId %d) new score %.1f and %.1f, previous score: %.1f and %.1f",
		event.LessonId, event.StudentId,
		unpackFloatRef(disciplineScore.Score.FirstScore), unpackFloatRef(disciplineScore.Score.SecondScore),
		unpackFloatRef(previousScore.FirstScore), unpackFloatRef(previousScore.SecondScore),
	)

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
		scoreChangesSendCount.Inc()
		err, newMessageId := handler.serviceContainer.ClientController.ScoreChangedAction(
			chatId, previousMessageIds[chatId], &disciplineScore, previousScore,
		)

		handler.debugLogger.Log(
			"ScoreChangedEventHandler: change (lessonId:%d, studentId %d) send update: message id %s, err: %v",
			event.LessonId, event.StudentId, newMessageId, err,
		)

		if err != nil {
			welcomeAuthorizedActionErrorCount.Inc()
			_, _ = fmt.Fprintf(handler.out, "ScoreChangedAction return error: %v\n", err)
		}

		if newMessageId != previousMessageIds[chatId] && (newMessageId != "" || err == nil) {
			handler.debugLogger.Log(
				"ScoreChangedEventHandler: change (lessonId:%d, studentId %d) save message id `%s` (err: %v)",
				event.LessonId, event.StudentId, newMessageId, err,
			)
			handler.scoreChangedMessageIdStorage.Set(event.StudentId, event.LessonId, chatId, newMessageId)
		}
	}

	handler.scoreChangedStateStorage.Set(event.StudentId, event.LessonId, newState)
}

func unpackFloatRef(input *float32) float32 {
	if input == nil {
		return 0
	}
	return *input
}
