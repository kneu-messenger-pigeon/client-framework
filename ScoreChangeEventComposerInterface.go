package framework

import (
	"github.com/kneu-messenger-pigeon/events"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
)

type ScoreChangeEventComposerInterface interface {
	Compose(event *events.ScoreChangedEvent, currentScore *scoreApi.Score) (previousScore *scoreApi.Score)
}
