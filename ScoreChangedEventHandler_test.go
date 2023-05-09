package framework

import (
	"bytes"
	"errors"
	"github.com/kneu-messenger-pigeon/client-framework/mocks"
	"github.com/kneu-messenger-pigeon/events"
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
	t.Run("success", func(t *testing.T) {
		clientController := mocks.NewClientControllerInterface(t)

		handler := ScoreChangedEventHandler{
			out: &bytes.Buffer{},
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		event := &events.ScoreChangedEvent{
			ScoreEvent: events.ScoreEvent{},
			Previous: struct {
				Value     float32
				IsAbsent  bool
				IsDeleted bool
			}{},
		}

		clientController.On("ScoreChangedAction", event).Return(nil)
		err := handler.Handle(event)

		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		clientController := mocks.NewClientControllerInterface(t)

		out := &bytes.Buffer{}
		handler := ScoreChangedEventHandler{
			out: out,
			serviceContainer: &ServiceContainer{
				ClientController: clientController,
			},
		}

		event := &events.ScoreChangedEvent{
			ScoreEvent: events.ScoreEvent{},
			Previous: struct {
				Value     float32
				IsAbsent  bool
				IsDeleted bool
			}{},
		}

		expectedError := errors.New("expected error")

		clientController.On("ScoreChangedAction", event).Return(expectedError)
		err := handler.Handle(event)

		// wait for async coroutine call
		time.Sleep(time.Millisecond * 40)

		assert.NoError(t, err)
		assert.Contains(t, out.String(), expectedError.Error())
	})
}

func TestScoreChangedEventHandler_Commit(t *testing.T) {
	handler := &ScoreChangedEventHandler{}
	assert.NoError(t, handler.Commit())
}
