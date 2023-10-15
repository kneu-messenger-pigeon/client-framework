package framework

import (
	"bytes"
	"errors"
	"github.com/kneu-messenger-pigeon/client-framework/mocks"
	"github.com/kneu-messenger-pigeon/client-framework/models"
	"github.com/kneu-messenger-pigeon/events"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUserAuthorizedEventHandler_GetExpectedMessageKey(t *testing.T) {
	handler := &UserAuthorizedEventHandler{}
	assert.Equal(t, events.UserAuthorizedEventName, handler.GetExpectedMessageKey())
}

func TestUserAuthorizedEventHandler_GetExpectedEventType(t *testing.T) {
	handler := &UserAuthorizedEventHandler{}
	assert.IsType(t, &events.UserAuthorizedEvent{}, handler.GetExpectedEventType())
}

func TestUserAuthorizedEventHandler_Handle(t *testing.T) {
	clientUserId := "test-client-id"

	t.Run("success_login", func(t *testing.T) {
		clientController := mocks.NewClientControllerInterface(t)
		userRepository := mocks.NewUserRepositoryInterface(t)

		out := &bytes.Buffer{}
		handler := UserAuthorizedEventHandler{
			out:        out,
			clientName: testClientName,
			repository: userRepository,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		expectedStudent := &models.Student{
			Id:         uint32(999),
			LastName:   "Потапенко",
			FirstName:  "Андрій",
			MiddleName: "Петрович",
			Gender:     models.Student_MALE,
		}

		event := &events.UserAuthorizedEvent{
			Client:       testClientName,
			ClientUserId: clientUserId,
			StudentId:    uint(expectedStudent.Id),
			LastName:     expectedStudent.LastName,
			FirstName:    expectedStudent.FirstName,
			MiddleName:   expectedStudent.MiddleName,
			Gender:       events.Gender(expectedStudent.Gender),
		}

		userRepository.On("SaveUser", clientUserId, expectedStudent).Return(nil, true)

		WelcomeAuthorizedActionError := errors.New("expected error")
		clientController.On("WelcomeAuthorizedAction", event).Return(WelcomeAuthorizedActionError)
		err := handler.Handle(event)

		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)

		assert.NoError(t, err)
		clientController.AssertNotCalled(t, "LogoutFinishedAction")
	})

	t.Run("success_logout", func(t *testing.T) {
		clientController := mocks.NewClientControllerInterface(t)
		userRepository := mocks.NewUserRepositoryInterface(t)

		out := &bytes.Buffer{}
		handler := UserAuthorizedEventHandler{
			out:        out,
			clientName: testClientName,
			repository: userRepository,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		expectedStudent := &models.Student{
			Id:         uint32(0),
			LastName:   "",
			FirstName:  "",
			MiddleName: "",
			Gender:     models.Student_UNKNOWN,
		}

		event := &events.UserAuthorizedEvent{
			Client:       testClientName,
			ClientUserId: clientUserId,
			StudentId:    uint(expectedStudent.Id),
			LastName:     expectedStudent.LastName,
			FirstName:    expectedStudent.FirstName,
			MiddleName:   expectedStudent.MiddleName,
			Gender:       events.Gender(expectedStudent.Gender),
		}

		userRepository.On("SaveUser", clientUserId, expectedStudent).Return(nil, true)

		expectedError := errors.New("expected error")
		clientController.On("LogoutFinishedAction", event).Return(expectedError)
		err := handler.Handle(event)

		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)

		assert.NoError(t, err)
		clientController.AssertNotCalled(t, "WelcomeAuthorizedAction")

		assert.Contains(t, out.String(), expectedError.Error())
	})

	t.Run("error", func(t *testing.T) {
		expectedError := errors.New("expected error")

		userRepository := mocks.NewUserRepositoryInterface(t)

		out := &bytes.Buffer{}
		handler := UserAuthorizedEventHandler{
			out:        out,
			clientName: testClientName,
			repository: userRepository,
		}

		expectedStudent := &models.Student{
			Id:         uint32(999),
			LastName:   "Потапенко",
			FirstName:  "Андрій",
			MiddleName: "Петрович",
			Gender:     models.Student_MALE,
		}

		event := &events.UserAuthorizedEvent{
			Client:       testClientName,
			ClientUserId: clientUserId,
			StudentId:    uint(expectedStudent.Id),
			LastName:     expectedStudent.LastName,
			FirstName:    expectedStudent.FirstName,
			MiddleName:   expectedStudent.MiddleName,
			Gender:       events.Gender(expectedStudent.Gender),
		}

		userRepository.On("SaveUser", clientUserId, expectedStudent).Return(expectedError, false)
		err := handler.Handle(event)

		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestUserAuthorizedEventHandler_Commit(t *testing.T) {
	userRepository := mocks.NewUserRepositoryInterface(t)
	handler := &UserAuthorizedEventHandler{
		repository: userRepository,
	}

	expectedError := errors.New("expected error")
	userRepository.On("Commit").Return(expectedError)

	actualError := handler.Commit()
	assert.Error(t, actualError)
	assert.Equal(t, expectedError, actualError)
}
