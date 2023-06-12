package framework

import (
	"bytes"
	"errors"
	"github.com/kneu-messenger-pigeon/client-framework/mocks"
	"github.com/kneu-messenger-pigeon/client-framework/models"
	"github.com/kneu-messenger-pigeon/events"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	"github.com/kneu-messenger-pigeon/score-client"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScoreChangedEventHandler_GetExpectedMessageKey(t *testing.T) {
	handler := &ScoreChangedEventHandler{}
	assert.Equal(t, events.ScoreChangedEventName, handler.GetExpectedMessageKey())
}

func TestScoreChangedEventHandler_GetExpectedEventType(t *testing.T) {
	handler := &ScoreChangedEventHandler{}
	assert.IsType(t, &events.ScoreChangedEvent{}, handler.GetExpectedEventType())
}

func TestScoreChangedEventHandler_Handle(t *testing.T) {
	t.Run("success_created_scores_no_previous_message", func(t *testing.T) {
		event := events.ScoreChangedEvent{
			ScoreEvent: events.ScoreEvent{
				Id:           112233,
				StudentId:    123,
				LessonId:     150,
				LessonPart:   1,
				DisciplineId: 234,
				Year:         2028,
				Semester:     1,
				ScoreValue: events.ScoreValue{
					Value:     2.5,
					IsAbsent:  false,
					IsDeleted: false,
				},
				UpdatedAt: time.Date(2028, time.Month(11), 18, 14, 30, 40, 0, time.Local),
				SyncedAt:  time.Date(2028, time.Month(11), 18, 14, 35, 13, 0, time.Local),
			},
			Previous: events.ScoreValue{
				Value:     0,
				IsAbsent:  false,
				IsDeleted: true,
			},
		}

		disciplineScore := scoreApi.DisciplineScore{
			Discipline: scoreApi.Discipline{
				Id:   int(event.DisciplineId),
				Name: "Капітал!",
			},
			Score: scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				FirstScore: floatPointer(2.5),
			},
		}

		chatIds := []string{
			"test-chat-id-1",
			"test-chat-id-2",
		}

		expectedMessageIds := []string{
			"chat-message-id-1",
			"chat-message-id-2",
		}

		previousScore := scoreApi.Score{}
		clientController := mocks.NewClientControllerInterface(t)

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return(chatIds)

		scoreClient := score.NewMockClientInterface(t)
		scoreClient.On("GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId)).
			Return(disciplineScore, nil)

		scoreChangeEventComposer := mocks.NewScoreChangeEventComposerInterface(t)
		scoreChangeEventComposer.On("Compose", &event, &disciplineScore.Score).Return(previousScore)

		messageIdStorage := mocks.NewScoreChangedMessageIdStorageInterface(t)
		messageIdStorage.On("GetAll", event.StudentId, event.LessonId).Return(models.ScoreChangedMessageMap{})

		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[0], expectedMessageIds[0]).Return()
		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[1], expectedMessageIds[1]).Return()

		handler := ScoreChangedEventHandler{
			out:                          &bytes.Buffer{},
			repository:                   userRepository,
			scoreClient:                  scoreClient,
			scoreChangedEventComposer:    scoreChangeEventComposer,
			scoreChangedMessageIdStorage: messageIdStorage,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		clientController.On("ScoreChangedAction", chatIds[0], "", &disciplineScore, &previousScore).
			Return(nil, expectedMessageIds[0])
		clientController.On("ScoreChangedAction", chatIds[1], "", &disciplineScore, &previousScore).
			Return(nil, expectedMessageIds[1])

		err := handler.Handle(&event)
		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)
		assert.NoError(t, err)
	})

	t.Run("success_created_scores_has_previous_message", func(t *testing.T) {
		event := events.ScoreChangedEvent{
			ScoreEvent: events.ScoreEvent{
				Id:           112233,
				StudentId:    123,
				LessonId:     150,
				LessonPart:   1,
				DisciplineId: 234,
				Year:         2028,
				Semester:     1,
				ScoreValue: events.ScoreValue{
					Value:     2.5,
					IsAbsent:  false,
					IsDeleted: false,
				},
				UpdatedAt: time.Date(2028, time.Month(11), 18, 14, 30, 40, 0, time.Local),
				SyncedAt:  time.Date(2028, time.Month(11), 18, 14, 35, 13, 0, time.Local),
			},
			Previous: events.ScoreValue{
				Value:     0,
				IsAbsent:  false,
				IsDeleted: true,
			},
		}
		disciplineScore := scoreApi.DisciplineScore{
			Discipline: scoreApi.Discipline{
				Id:   int(event.DisciplineId),
				Name: "Капітал!",
			},
			Score: scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				FirstScore: floatPointer(2.5),
			},
		}

		chatIds := []string{
			"test-chat-id-1",
			"test-chat-id-2",
		}

		expectedMessageIds := []string{
			"chat-message-id-1",
			"chat-message-id-2",
		}

		previousScore := scoreApi.Score{}

		clientController := mocks.NewClientControllerInterface(t)

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return(chatIds)

		scoreClient := score.NewMockClientInterface(t)
		scoreClient.On("GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId)).Return(disciplineScore, nil)

		scoreChangeEventComposer := mocks.NewScoreChangeEventComposerInterface(t)
		scoreChangeEventComposer.On("Compose", &event, &disciplineScore.Score).Return(previousScore)

		messageIdStorage := mocks.NewScoreChangedMessageIdStorageInterface(t)
		messageIdStorage.On("GetAll", event.StudentId, event.LessonId).
			Return(models.ScoreChangedMessageMap{
				chatIds[0]: expectedMessageIds[0],
				chatIds[1]: expectedMessageIds[1],
			})

		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[0], expectedMessageIds[0]).Return()
		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[1], expectedMessageIds[1]).Return()

		handler := ScoreChangedEventHandler{
			out:                          &bytes.Buffer{},
			repository:                   userRepository,
			scoreClient:                  scoreClient,
			scoreChangedEventComposer:    scoreChangeEventComposer,
			scoreChangedMessageIdStorage: messageIdStorage,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		clientController.On("ScoreChangedAction", chatIds[0], expectedMessageIds[0], &disciplineScore, &previousScore).
			Return(nil, expectedMessageIds[0])
		clientController.On("ScoreChangedAction", chatIds[1], expectedMessageIds[1], &disciplineScore, &previousScore).
			Return(nil, expectedMessageIds[1])

		err := handler.Handle(&event)
		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)
		assert.NoError(t, err)
	})

	t.Run("success_created_scores_remove_previous_message", func(t *testing.T) {
		event := events.ScoreChangedEvent{
			ScoreEvent: events.ScoreEvent{
				Id:           112233,
				StudentId:    123,
				LessonId:     150,
				LessonPart:   1,
				DisciplineId: 234,
				Year:         2028,
				Semester:     1,
				ScoreValue: events.ScoreValue{
					Value:     2.5,
					IsAbsent:  false,
					IsDeleted: false,
				},
				UpdatedAt: time.Date(2028, time.Month(11), 18, 14, 30, 40, 0, time.Local),
				SyncedAt:  time.Date(2028, time.Month(11), 18, 14, 35, 13, 0, time.Local),
			},
			Previous: events.ScoreValue{
				Value:     0,
				IsAbsent:  false,
				IsDeleted: true,
			},
		}
		disciplineScore := scoreApi.DisciplineScore{
			Discipline: scoreApi.Discipline{
				Id:   int(event.DisciplineId),
				Name: "Капітал!",
			},
			Score: scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				FirstScore: floatPointer(2.5),
			},
		}

		chatIds := []string{
			"test-chat-id-1",
			"test-chat-id-2",
		}

		previousMessageIds := []string{
			"chat-message-id-1",
			"chat-message-id-2",
		}

		previousScore := scoreApi.Score{}

		clientController := mocks.NewClientControllerInterface(t)

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return(chatIds)

		scoreClient := score.NewMockClientInterface(t)
		scoreClient.On("GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId)).Return(disciplineScore, nil)

		scoreChangeEventComposer := mocks.NewScoreChangeEventComposerInterface(t)
		scoreChangeEventComposer.On("Compose", &event, &disciplineScore.Score).Return(previousScore)

		messageIdStorage := mocks.NewScoreChangedMessageIdStorageInterface(t)
		messageIdStorage.On("GetAll", event.StudentId, event.LessonId).
			Return(models.ScoreChangedMessageMap{
				chatIds[0]: previousMessageIds[0],
				chatIds[1]: previousMessageIds[1],
			})

		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[0], "").Return()
		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[1], "").Return()

		handler := ScoreChangedEventHandler{
			out:                          &bytes.Buffer{},
			repository:                   userRepository,
			scoreClient:                  scoreClient,
			scoreChangedEventComposer:    scoreChangeEventComposer,
			scoreChangedMessageIdStorage: messageIdStorage,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		clientController.On("ScoreChangedAction", chatIds[0], previousMessageIds[0], &disciplineScore, &previousScore).
			Return(nil, "")
		clientController.On("ScoreChangedAction", chatIds[1], previousMessageIds[1], &disciplineScore, &previousScore).
			Return(nil, "")

		err := handler.Handle(&event)
		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)
		assert.NoError(t, err)
	})

	t.Run("success_created_scores_controller_action_err", func(t *testing.T) {
		expectedError := errors.New("test expected error")

		out := &bytes.Buffer{}

		event := events.ScoreChangedEvent{
			ScoreEvent: events.ScoreEvent{
				Id:           112233,
				StudentId:    123,
				LessonId:     150,
				LessonPart:   1,
				DisciplineId: 234,
				Year:         2028,
				Semester:     1,
				ScoreValue: events.ScoreValue{
					Value:     2.5,
					IsAbsent:  false,
					IsDeleted: false,
				},
				UpdatedAt: time.Date(2028, time.Month(11), 18, 14, 30, 40, 0, time.Local),
				SyncedAt:  time.Date(2028, time.Month(11), 18, 14, 35, 13, 0, time.Local),
			},
			Previous: events.ScoreValue{
				Value:     0,
				IsAbsent:  false,
				IsDeleted: true,
			},
		}

		disciplineScore := scoreApi.DisciplineScore{
			Discipline: scoreApi.Discipline{
				Id:   int(event.DisciplineId),
				Name: "Капітал!",
			},
			Score: scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				FirstScore: floatPointer(2.5),
			},
		}

		chatIds := []string{
			"test-chat-id-1",
			"test-chat-id-2",
		}

		expectedMessageIds := []string{
			"chat-message-id-1",
			"chat-message-id-2",
		}

		previousScore := scoreApi.Score{}
		clientController := mocks.NewClientControllerInterface(t)

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return(chatIds)

		scoreClient := score.NewMockClientInterface(t)
		scoreClient.On(
			"GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId),
		).Return(disciplineScore, nil)

		scoreChangeEventComposer := mocks.NewScoreChangeEventComposerInterface(t)
		scoreChangeEventComposer.On("Compose", &event, &disciplineScore.Score).Return(previousScore)

		messageIdStorage := mocks.NewScoreChangedMessageIdStorageInterface(t)
		messageIdStorage.On("GetAll", event.StudentId, event.LessonId).Return(models.ScoreChangedMessageMap{})
		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[0], expectedMessageIds[0]).Return()
		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[1], expectedMessageIds[1]).Return()

		handler := ScoreChangedEventHandler{
			out:                          out,
			repository:                   userRepository,
			scoreClient:                  scoreClient,
			scoreChangedEventComposer:    scoreChangeEventComposer,
			scoreChangedMessageIdStorage: messageIdStorage,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		clientController.On("ScoreChangedAction", chatIds[0], "", &disciplineScore, &previousScore).
			Return(expectedError, expectedMessageIds[0])
		clientController.On("ScoreChangedAction", chatIds[1], "", &disciplineScore, &previousScore).
			Return(nil, expectedMessageIds[1])

		err := handler.Handle(&event)
		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)
		assert.NoError(t, err)

		assert.Contains(t, out.String(), expectedError.Error())
	})

	t.Run("success_score_client_err", func(t *testing.T) {
		expectedError := errors.New("test expected error")

		event := &events.ScoreChangedEvent{
			ScoreEvent: events.ScoreEvent{
				Id:           112233,
				StudentId:    456,
				LessonId:     150,
				LessonPart:   1,
				DisciplineId: 234,
			},
		}

		clientController := mocks.NewClientControllerInterface(t)

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return([]string{"test-chat-id-1"})

		scoreClient := score.NewMockClientInterface(t)
		scoreClient.On(
			"GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId),
		).Return(scoreApi.DisciplineScore{}, expectedError)

		handler := ScoreChangedEventHandler{
			out:         &bytes.Buffer{},
			repository:  userRepository,
			scoreClient: scoreClient,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}
		actualErr := handler.Handle(event)
		assert.Error(t, actualErr)
		assert.Equal(t, expectedError, actualErr)
	})

	t.Run("success_no_chat_ids", func(t *testing.T) {
		event := &events.ScoreChangedEvent{
			ScoreEvent: events.ScoreEvent{
				StudentId: 456,
			},
		}

		clientController := mocks.NewClientControllerInterface(t)

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return([]string{})

		handler := ScoreChangedEventHandler{
			out:        &bytes.Buffer{},
			repository: userRepository,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}
		err := handler.Handle(event)
		assert.NoError(t, err)
	})

	t.Run("success_no_service_container", func(t *testing.T) {
		handler := ScoreChangedEventHandler{
			out:              &bytes.Buffer{},
			serviceContainer: &ServiceContainer{},
		}
		err := handler.Handle(&events.ScoreChangedEvent{})
		assert.NoError(t, err)
	})
}

func TestScoreChangedEventHandler_Commit(t *testing.T) {
	handler := &ScoreChangedEventHandler{}
	assert.NoError(t, handler.Commit())
}
