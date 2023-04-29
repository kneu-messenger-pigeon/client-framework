package framework

import "github.com/kneu-messenger-pigeon/events"

type ClientControllerInterface interface {
	ExecutableInterface
	ScoreChangedAction(event *events.ScoreChangedEvent) error
	UserAuthorizedAction(event *events.UserAuthorizedEvent) error
}
