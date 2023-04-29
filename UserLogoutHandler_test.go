package framework

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/kneu-messenger-pigeon/events"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestUserLogoutHandler_Handle(t *testing.T) {
	var matchContext = mock.MatchedBy(func(ctx context.Context) bool { return true })

	messageMatcher := func(expectedMessage events.UserAuthorizedEvent) interface{} {
		return mock.MatchedBy(func(message kafka.Message) bool {
			actualMessage := events.UserAuthorizedEvent{}
			err := json.Unmarshal(message.Value, &actualMessage)
			return assert.Equal(t, events.UserAuthorizedEventName, string(message.Key)) &&
				assert.NoErrorf(t, err, "Failed to parse as UserAuthorizedEvent: %v", message) &&
				assert.Equal(t, expectedMessage, actualMessage)
		})
	}

	t.Run("success", func(t *testing.T) {
		expectedMessage := events.UserAuthorizedEvent{
			Client:       "test-client",
			ClientUserId: "test-client-id",
			StudentId:    0,
			LastName:     "",
			FirstName:    "",
			MiddleName:   "",
			Gender:       events.UnknownGender,
		}

		writer := events.NewMockWriterInterface(t)
		writer.On("WriteMessages", matchContext, messageMatcher(expectedMessage)).Return(nil)

		handler := UserLogoutHandler{
			out:    &bytes.Buffer{},
			Client: "test-client",
			writer: writer,
		}

		err := handler.Handle(expectedMessage.ClientUserId)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		expectedMessage := events.UserAuthorizedEvent{
			Client:       "test-client",
			ClientUserId: "test-client-id",
			StudentId:    0,
			LastName:     "",
			FirstName:    "",
			MiddleName:   "",
			Gender:       events.UnknownGender,
		}
		expectedError := errors.New("test expected error")

		writer := events.NewMockWriterInterface(t)
		writer.On("WriteMessages", matchContext, messageMatcher(expectedMessage)).Return(expectedError)

		handler := UserLogoutHandler{
			out:    &bytes.Buffer{},
			Client: "test-client",
			writer: writer,
		}

		err := handler.Handle(expectedMessage.ClientUserId)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}
