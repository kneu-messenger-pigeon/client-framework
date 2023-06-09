package framework

import (
	"context"
	"encoding/json"
	"github.com/kneu-messenger-pigeon/events"
	"github.com/segmentio/kafka-go"
	"io"
)

type UserLogoutHandler struct {
	out    io.Writer
	Client string
	writer events.WriterInterface
}

func (handler UserLogoutHandler) Handle(clientUserId string) error {
	event := events.UserAuthorizedEvent{
		Client:       handler.Client,
		ClientUserId: clientUserId,
		StudentId:    0,
		LastName:     "",
		FirstName:    "",
		MiddleName:   "",
		Gender:       events.UnknownGender,
	}

	payload, _ := json.Marshal(event)
	return handler.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(events.UserAuthorizedEventName),
			Value: payload,
		},
	)
}
