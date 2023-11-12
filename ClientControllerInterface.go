package framework

import (
	delayedDeleterContracts "github.com/kneu-messenger-pigeon/client-framework/delayedDeleter/contracts"
	"github.com/kneu-messenger-pigeon/events"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
)

type ClientControllerInterface interface {
	ExecutableInterface
	delayedDeleterContracts.DeleteHandlerInterface

	ScoreChangedAction(chatId string, previousMessageId string, disciplineScore *scoreApi.DisciplineScore, previousScore *scoreApi.Score) (err error, messageId string)
	WelcomeAuthorizedAction(event *events.UserAuthorizedEvent) error
	LogoutFinishedAction(event *events.UserAuthorizedEvent) error
}
