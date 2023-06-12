package framework

import "github.com/kneu-messenger-pigeon/client-framework/models"

type MessageComposerInterface interface {
	SetPostFilter(filter func(string) string)
	ComposeWelcomeAnonymousMessage(authUrl string) (error, string)
	ComposeWelcomeAuthorizedMessage(messageData models.UserAuthorizedMessageData) (error, string)
	ComposeDisciplinesListMessage(messageData models.DisciplinesListMessageData) (error, string)
	ComposeDisciplineScoresMessage(messageData models.DisciplinesScoresMessageData) (error, string)
	ComposeScoreChanged(messageData models.ScoreChangedMessageData) (error, string)
	ComposeLogoutFinishedMessage() (error, string)
}
