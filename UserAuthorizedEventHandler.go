package framework

import (
	"fmt"
	"github.com/kneu-messenger-pigeon/client-framework/models"
	"github.com/kneu-messenger-pigeon/events"
	"io"
)

type UserAuthorizedEventHandler struct {
	clientName       string
	repository       UserRepositoryInterface
	out              io.Writer
	serviceContainer *ServiceContainer
}

func (handler *UserAuthorizedEventHandler) GetExpectedMessageKey() string {
	return events.UserAuthorizedEventName
}

func (handler *UserAuthorizedEventHandler) GetExpectedEventType() any {
	return &events.UserAuthorizedEvent{}
}

func (handler *UserAuthorizedEventHandler) Commit() error {
	return handler.repository.Commit()
}

func (handler *UserAuthorizedEventHandler) Handle(s any) (err error) {
	event := s.(*events.UserAuthorizedEvent)
	if event.Client == handler.clientName {
		var hasChanges bool
		err, hasChanges = handler.repository.SaveUser(event.ClientUserId, &models.Student{
			Id:         uint32(event.StudentId),
			LastName:   event.LastName,
			FirstName:  event.FirstName,
			MiddleName: event.MiddleName,
			Gender:     models.Student_GenderType(event.Gender),
		})

		if err == nil && handler.serviceContainer != nil && handler.serviceContainer.ClientController != nil {
			go handler.callControllerAction(event)
		}

		if err == nil && hasChanges {
			if event.StudentId != 0 {
				userCount.Inc()
				loginCount.Inc()
			} else {
				logoutCount.Inc()
				userCount.Dec()
			}
		}
	}

	return err
}

func (handler *UserAuthorizedEventHandler) callControllerAction(event *events.UserAuthorizedEvent) {
	var err error
	if event.StudentId != 0 {
		err = handler.serviceContainer.ClientController.WelcomeAuthorizedAction(event)
		if err != nil {
			welcomeAuthorizedActionErrorCount.Inc()
		}
	} else {
		err = handler.serviceContainer.ClientController.LogoutFinishedAction(event)
	}

	if err != nil {
		_, _ = fmt.Fprintf(handler.out, "UserAuthorizedAction return error: %v\n", err)
	}
}
