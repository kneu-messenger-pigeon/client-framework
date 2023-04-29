package framework

import "github.com/kneu-messenger-pigeon/events"

type ClientControllerInterface interface {
	ExecutableInterface
	ScoreChangedAction(event *events.ScoreChangedEvent) error
	WelcomeAuthorizedAction(event *events.UserAuthorizedEvent) error
	LogoutFinishedAction(event *events.UserAuthorizedEvent) error
}
