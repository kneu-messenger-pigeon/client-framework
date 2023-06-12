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
	repository                   UserRepositoryInterface
	scoreClient                  score.ClientInterface
	scoreChangedEventComposer    ScoreChangeEventComposerInterface
	scoreChangedMessageIdStorage ScoreChangedMessageIdStorageInterface
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
	if len(chatIds) == 0 {
		return nil
	}

	/**
	 * @todo
	 * - implement storing send message chat id and restore it from storage
	 */
	disciplineScore, err := handler.scoreClient.GetStudentScore(
		uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId),
	)
	if err != nil {
		return err
	}

	previousScore := handler.scoreChangedEventComposer.Compose(event, &disciplineScore.Score)

	previousMessageIds := handler.scoreChangedMessageIdStorage.GetAll(event.StudentId, event.LessonId)

	for _, chatId := range chatIds {
		go handler.callControllerAction(
			event.StudentId, chatId, previousMessageIds[chatId],
			&disciplineScore, &previousScore,
		)
	}

	return nil
}

func (handler *ScoreChangedEventHandler) callControllerAction(
	studentId uint, chatId string, previousMessageId string,
	score *scoreApi.DisciplineScore, previousScore *scoreApi.Score,
) {
	err, messageId := handler.serviceContainer.ClientController.ScoreChangedAction(
		chatId, previousMessageId, score, previousScore,
	)

	if err != nil {
		_, _ = fmt.Fprintf(handler.out, "ScoreChangedAction return error: %v", err)
	}

	if messageId != "" || previousMessageId != "" {
		handler.scoreChangedMessageIdStorage.Set(studentId, uint(score.Score.Lesson.Id), chatId, messageId)
	}
}
