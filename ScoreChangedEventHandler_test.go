package framework

import (
	"bytes"
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/kneu-messenger-pigeon/client-framework/mocks"
	"github.com/kneu-messenger-pigeon/client-framework/models"
	"github.com/kneu-messenger-pigeon/events"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	scoreMocks "github.com/kneu-messenger-pigeon/score-client/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"runtime"
	"strings"
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

		previousScore := &scoreApi.Score{}
		clientController := mocks.NewClientControllerInterface(t)

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return(chatIds)

		scoreClient := scoreMocks.NewClientInterface(t)
		scoreClient.On("GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId)).
			Return(disciplineScore, nil)

		scoreChangeEventComposer := mocks.NewScoreChangeEventComposerInterface(t)
		scoreChangeEventComposer.On("SavePreviousScore", &event).Return()
		scoreChangeEventComposer.On("BothPreviousScoresSaved", &event).Return(false)
		scoreChangeEventComposer.On("Compose", &event, &disciplineScore.Score).Return(previousScore)

		stateStorage := mocks.NewScoreChangedStateStorageInterface(t)
		stateStorage.On("Get", event.StudentId, event.LessonId).Once().Return("")
		stateStorage.On("Set", event.StudentId, event.LessonId, mock.Anything).Once().Return()

		messageIdStorage := mocks.NewScoreChangedMessageIdStorageInterface(t)
		messageIdStorage.On("GetAll", event.StudentId, event.LessonId).Return(models.ScoreChangedMessageMap{})

		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[0], expectedMessageIds[0]).Return()
		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[1], expectedMessageIds[1]).Return()

		handler := ScoreChangedEventHandler{
			out:                          &bytes.Buffer{},
			debugLogger:                  &DebugLogger{},
			repository:                   userRepository,
			scoreClient:                  scoreClient,
			scoreChangedEventComposer:    scoreChangeEventComposer,
			scoreChangedStateStorage:     stateStorage,
			scoreChangedMessageIdStorage: messageIdStorage,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
			waitingForAnotherScoreTime: time.Millisecond * 100,
		}

		clientController.On("ScoreChangedAction", chatIds[0], "", &disciplineScore, previousScore).
			Return(nil, expectedMessageIds[0])
		clientController.On("ScoreChangedAction", chatIds[1], "", &disciplineScore, previousScore).
			Return(nil, expectedMessageIds[1])

		err := handler.Handle(&event)
		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)
		clientController.AssertNotCalled(t, "ScoreChangedAction", chatIds[0], "", &disciplineScore, previousScore)
		time.Sleep(handler.waitingForAnotherScoreTime)
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

		existsMessageIds := []string{
			"chat-message-id-1",
			"chat-message-id-2",
		}

		previousScore := &scoreApi.Score{}

		clientController := mocks.NewClientControllerInterface(t)

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return(chatIds)

		scoreClient := scoreMocks.NewClientInterface(t)
		scoreClient.On("GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId)).Return(disciplineScore, nil)

		scoreChangeEventComposer := mocks.NewScoreChangeEventComposerInterface(t)
		scoreChangeEventComposer.On("SavePreviousScore", &event).Return()
		scoreChangeEventComposer.On("BothPreviousScoresSaved", &event).Return(false)
		scoreChangeEventComposer.On("Compose", &event, &disciplineScore.Score).Return(previousScore)

		stateStorage := mocks.NewScoreChangedStateStorageInterface(t)
		stateStorage.On("Get", event.StudentId, event.LessonId).Once().Return("123")
		stateStorage.On("Set", event.StudentId, event.LessonId, mock.Anything).Once().Return()

		messageIdStorage := mocks.NewScoreChangedMessageIdStorageInterface(t)
		messageIdStorage.On("GetAll", event.StudentId, event.LessonId).
			Return(models.ScoreChangedMessageMap{
				chatIds[0]: existsMessageIds[0],
				chatIds[1]: existsMessageIds[1],
			})

		handler := ScoreChangedEventHandler{
			out:                          &bytes.Buffer{},
			debugLogger:                  &DebugLogger{},
			repository:                   userRepository,
			scoreClient:                  scoreClient,
			scoreChangedEventComposer:    scoreChangeEventComposer,
			scoreChangedStateStorage:     stateStorage,
			scoreChangedMessageIdStorage: messageIdStorage,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		clientController.On("ScoreChangedAction", chatIds[0], existsMessageIds[0], &disciplineScore, previousScore).
			Return(nil, existsMessageIds[0])
		clientController.On("ScoreChangedAction", chatIds[1], existsMessageIds[1], &disciplineScore, previousScore).
			Return(nil, existsMessageIds[1])

		err := handler.Handle(&event)
		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)
		assert.NoError(t, err)

		// no changes - no need to save in storage
		messageIdStorage.AssertNotCalled(t, "Set")
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

		previousScore := &scoreApi.Score{}

		clientController := mocks.NewClientControllerInterface(t)

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return(chatIds)

		scoreClient := scoreMocks.NewClientInterface(t)
		scoreClient.On("GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId)).Return(disciplineScore, nil)

		scoreChangeEventComposer := mocks.NewScoreChangeEventComposerInterface(t)
		scoreChangeEventComposer.On("SavePreviousScore", &event).Return()
		scoreChangeEventComposer.On("BothPreviousScoresSaved", &event).Return(false)
		scoreChangeEventComposer.On("Compose", &event, &disciplineScore.Score).Return(previousScore)

		stateStorage := mocks.NewScoreChangedStateStorageInterface(t)
		stateStorage.On("Get", event.StudentId, event.LessonId).Once().Return("123")
		stateStorage.On("Set", event.StudentId, event.LessonId, mock.Anything).Once().Return()

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
			debugLogger:                  &DebugLogger{},
			repository:                   userRepository,
			scoreClient:                  scoreClient,
			scoreChangedEventComposer:    scoreChangeEventComposer,
			scoreChangedStateStorage:     stateStorage,
			scoreChangedMessageIdStorage: messageIdStorage,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		clientController.On("ScoreChangedAction", chatIds[0], previousMessageIds[0], &disciplineScore, previousScore).
			Return(nil, "")
		clientController.On("ScoreChangedAction", chatIds[1], previousMessageIds[1], &disciplineScore, previousScore).
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
			"test-chat-id-3",
		}

		expectedMessageIds := []string{
			"chat-message-id-1",
			"chat-message-id-2",
			"chat-message-id-3",
		}

		previousScore := &scoreApi.Score{}
		clientController := mocks.NewClientControllerInterface(t)

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return(chatIds)

		scoreClient := scoreMocks.NewClientInterface(t)
		scoreClient.On(
			"GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId),
		).Return(disciplineScore, nil)

		scoreChangeEventComposer := mocks.NewScoreChangeEventComposerInterface(t)
		scoreChangeEventComposer.On("SavePreviousScore", &event).Return()
		scoreChangeEventComposer.On("BothPreviousScoresSaved", &event).Return(false)
		scoreChangeEventComposer.On("Compose", &event, &disciplineScore.Score).Return(previousScore)

		stateStorage := mocks.NewScoreChangedStateStorageInterface(t)
		stateStorage.On("Get", event.StudentId, event.LessonId).Once().Return("123")
		stateStorage.On("Set", event.StudentId, event.LessonId, mock.Anything).Once().Return()

		messageIdStorage := mocks.NewScoreChangedMessageIdStorageInterface(t)
		messageIdStorage.On("GetAll", event.StudentId, event.LessonId).Return(models.ScoreChangedMessageMap{})
		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[0], expectedMessageIds[0]).Return()
		messageIdStorage.On("Set", event.StudentId, event.LessonId, chatIds[1], expectedMessageIds[1]).Return()

		clientController.On("ScoreChangedAction", chatIds[0], "", &disciplineScore, previousScore).
			Return(expectedError, expectedMessageIds[0])
		clientController.On("ScoreChangedAction", chatIds[1], "", &disciplineScore, previousScore).
			Return(nil, expectedMessageIds[1])

		clientController.On("ScoreChangedAction", chatIds[2], "", &disciplineScore, previousScore).
			Return(expectedError, "")

		handler := ScoreChangedEventHandler{
			out:                          out,
			debugLogger:                  &DebugLogger{},
			repository:                   userRepository,
			scoreClient:                  scoreClient,
			scoreChangedEventComposer:    scoreChangeEventComposer,
			scoreChangedStateStorage:     stateStorage,
			scoreChangedMessageIdStorage: messageIdStorage,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
			waitingForAnotherScoreTime: time.Millisecond * 500,
		}

		timeToCall := time.After(time.Millisecond * 550)
		err := handler.Handle(&event)
		assert.NoError(t, err)
		// wait for async coroutine call
		time.Sleep(time.Millisecond * 200)
		runtime.Gosched()
		clientController.AssertNotCalled(t, chatIds[0], "", &disciplineScore, previousScore)
		<-timeToCall
		runtime.Gosched()

		assert.Contains(t, out.String(), expectedError.Error())
		assert.Equal(t, 2, strings.Count(out.String(), expectedError.Error()))
		messageIdStorage.AssertNotCalled(t, "Set", event.StudentId, event.LessonId, chatIds[2], expectedMessageIds[2])
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

		scoreClient := scoreMocks.NewClientInterface(t)
		scoreClient.On(
			"GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId),
		).Return(scoreApi.DisciplineScore{}, expectedError)

		scoreChangeEventComposer := mocks.NewScoreChangeEventComposerInterface(t)
		scoreChangeEventComposer.On("SavePreviousScore", event).Return()
		scoreChangeEventComposer.On("BothPreviousScoresSaved", event).Return(false)

		out := &bytes.Buffer{}

		handler := ScoreChangedEventHandler{
			out:                       out,
			debugLogger:               &DebugLogger{},
			repository:                userRepository,
			scoreClient:               scoreClient,
			scoreChangedEventComposer: scoreChangeEventComposer,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}
		actualErr := handler.Handle(event)
		assert.NoError(t, actualErr)
		runtime.Gosched()
		time.Sleep(time.Millisecond * 40)
		assert.Contains(t, out.String(), "GetStudentScore return error: "+expectedError.Error())
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
			out:         &bytes.Buffer{},
			debugLogger: &DebugLogger{},
			repository:  userRepository,
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
			debugLogger:      &DebugLogger{},
			serviceContainer: &ServiceContainer{},
		}
		err := handler.Handle(&events.ScoreChangedEvent{})
		assert.NoError(t, err)
	})

	t.Run("same_change_state_no_send_message", func(t *testing.T) {
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

		previousScore := &scoreApi.Score{
			Lesson: disciplineScore.Score.Lesson,
		}

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event.StudentId).Return([]string{"test-chat-id-1"})

		scoreClient := scoreMocks.NewClientInterface(t)
		scoreClient.
			On("GetStudentScore", uint32(event.StudentId), int(event.DisciplineId), int(event.LessonId)).
			Once().
			Return(disciplineScore, nil)

		scoreChangeEventComposer := mocks.NewScoreChangeEventComposerInterface(t)
		scoreChangeEventComposer.On("SavePreviousScore", &event).Return()
		scoreChangeEventComposer.On("BothPreviousScoresSaved", &event).Return(false)
		scoreChangeEventComposer.On("Compose", &event, &disciplineScore.Score).Return(previousScore)

		previousState := CalculateState(&disciplineScore.Score, previousScore)

		stateStorage := mocks.NewScoreChangedStateStorageInterface(t)
		stateStorage.On("Get", event.StudentId, event.LessonId).Once().Return(previousState)

		clientController := mocks.NewClientControllerInterface(t)

		handler := ScoreChangedEventHandler{
			out:                          &bytes.Buffer{},
			debugLogger:                  &DebugLogger{},
			repository:                   userRepository,
			scoreClient:                  scoreClient,
			scoreChangedEventComposer:    scoreChangeEventComposer,
			scoreChangedStateStorage:     stateStorage,
			scoreChangedMessageIdStorage: mocks.NewScoreChangedMessageIdStorageInterface(t),
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		err := handler.Handle(&event)
		assert.NoError(t, err)
		runtime.Gosched()
		time.Sleep(time.Millisecond * 40)
		clientController.AssertNotCalled(t, "ScoreChangedAction")
	})

	t.Run("success_created_scores_race_condition", func(t *testing.T) {
		event1 := events.ScoreChangedEvent{
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

		event2 := events.ScoreChangedEvent{
			ScoreEvent: events.ScoreEvent{
				Id:           event1.Id,
				Year:         event1.Year,
				Semester:     event1.Semester,
				StudentId:    event1.StudentId,
				LessonId:     event1.LessonId,
				DisciplineId: event1.DisciplineId,
				LessonPart:   2,
				ScoreValue: events.ScoreValue{
					Value:     3,
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

		disciplineScore1 := scoreApi.DisciplineScore{
			Discipline: scoreApi.Discipline{
				Id:   int(event1.DisciplineId),
				Name: "Капітал!",
			},
			Score: scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event1.LessonId),
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

		disciplineScore2 := scoreApi.DisciplineScore{
			Discipline: scoreApi.Discipline{
				Id:   int(event1.DisciplineId),
				Name: "Капітал!",
			},
			Score: scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   int(event1.LessonId),
					Date: time.Date(2023, time.Month(2), 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        5,
						ShortName: "МК",
						LongName:  "Модульний контроль.",
					},
				},
				SecondScore: floatPointer(3),
			},
		}

		chatIds := []string{
			"test-chat-id-1",
			"test-chat-id-2",
		}

		createdMessageIds := []string{
			"chat-message-id-1",
			"chat-message-id-2",
		}

		previousScore := scoreApi.Score{
			Lesson: disciplineScore1.Score.Lesson,
		}

		userRepository := mocks.NewUserRepositoryInterface(t)
		userRepository.On("GetClientUserIds", event1.StudentId).Return(chatIds)

		scoreClient := scoreMocks.NewClientInterface(t)
		scoreClient.
			On("GetStudentScore", uint32(event1.StudentId), int(event1.DisciplineId), int(event1.LessonId)).
			Once().
			Return(disciplineScore1, nil)

		redisClient := redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    miniredis.RunT(t).Addr(),
		})

		firstMessageSendingWait := make(chan time.Time)

		clientController := mocks.NewClientControllerInterface(t)
		clientController.On("ScoreChangedAction", chatIds[0], "", &disciplineScore1, &previousScore).
			Once().WaitUntil(firstMessageSendingWait).
			Return(nil, createdMessageIds[0])

		clientController.On("ScoreChangedAction", chatIds[1], "", &disciplineScore1, &previousScore).
			Once().WaitUntil(firstMessageSendingWait).
			Return(nil, createdMessageIds[1])

		handler := ScoreChangedEventHandler{
			out:         &bytes.Buffer{},
			debugLogger: &DebugLogger{},
			repository:  userRepository,
			scoreClient: scoreClient,
			scoreChangedEventComposer: &ScoreChangeEventComposer{
				out:           &bytes.Buffer{},
				redis:         redisClient,
				storageExpire: time.Minute * 5,
			},
			scoreChangedStateStorage: &ScoreChangedStateStorage{
				redis:         redisClient,
				storageExpire: time.Minute,
			},
			scoreChangedMessageIdStorage: &ScoreChangedMessageIdStorage{
				out:           &bytes.Buffer{},
				redis:         redisClient,
				storageExpire: time.Minute,
			},
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		err := handler.Handle(&event1)
		assert.NoError(t, err)
		runtime.Gosched()
		time.Sleep(time.Millisecond * 40)
		clientController.AssertNumberOfCalls(t, "ScoreChangedAction", 1)

		scoreClient.On("GetStudentScore", uint32(event2.StudentId), int(event2.DisciplineId), int(event2.LessonId)).
			Once().
			Return(disciplineScore2, nil)

		// process seconds message in queue
		clientController.On("ScoreChangedAction", chatIds[0], createdMessageIds[0], &disciplineScore2, &previousScore).
			Once().
			Return(nil, createdMessageIds[0])

		clientController.On("ScoreChangedAction", chatIds[1], createdMessageIds[1], &disciplineScore2, &previousScore).
			Once().
			Return(nil, createdMessageIds[1])

		err = handler.Handle(&event2)
		runtime.Gosched() // wait for async coroutine call
		time.Sleep(time.Millisecond * 100)

		// ensure that handler for second iteration message was not called yet (while first iteration isn't finished)
		clientController.AssertNotCalled(t, "ScoreChangedAction", chatIds[0], createdMessageIds[0], &disciplineScore2, &previousScore)
		clientController.AssertNotCalled(t, "ScoreChangedAction", chatIds[1], createdMessageIds[1], &disciplineScore2, &previousScore)

		// finish handlers from first iteration and expect it will run handlers from second iteration.
		firstMessageSendingWait <- time.Time{}
		firstMessageSendingWait <- time.Time{}

		runtime.Gosched() // wait for async coroutine call
		time.Sleep(time.Millisecond * 100)

		assert.NoError(t, err)
	})
}

func TestScoreChangedEventHandler_Commit(t *testing.T) {
	handler := &ScoreChangedEventHandler{}
	assert.NoError(t, handler.Commit())
}
