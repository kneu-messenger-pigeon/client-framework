package framework

import (
	"bytes"
	"embed"
	_ "embed"
	"github.com/kneu-messenger-pigeon/client-framework/models"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	"strconv"
	"text/template"
	"time"
)

//go:embed templates/*.md
var templates embed.FS

type MessageComposer struct {
	templates  *template.Template
	postFilter func(string) string
}

type MessageComposerConfig struct {
}

func NewMessageComposer(config MessageComposerConfig) *MessageComposer {
	composer := &MessageComposer{}

	composer.templates = template.Must(
		template.New("").
			Funcs(composer.getFunctionMap()).
			ParseFS(templates, "templates/*.md"),
	)

	composer.postFilter = func(i string) string {
		return i
	}

	/*
		for _, tpl := range composer.templates.Templates() {
			templateString := tpl.Tree.Root.String()

			template.Must(tpl.Parse(templateString))
			template.Must(composer.templates.AddParseTree(tpl.Name(), tpl.Tree))
		}

	*/

	return composer
}

func (composer *MessageComposer) SetPostFilter(filter func(string) string) {
	composer.postFilter = filter
}

func (composer *MessageComposer) ComposeWelcomeAnonymousMessage(authUrl string) (error, string) {
	return composer.compose("WelcomeAnonymous.md", authUrl)
}

func (composer *MessageComposer) ComposeWelcomeAuthorizedMessage(messageData models.UserAuthorizedMessageData) (error, string) {
	return composer.compose("WelcomeAuthorized.md", messageData)
}

func (composer *MessageComposer) ComposeDisciplinesListMessage(messageData models.DisciplinesListMessageData) (error, string) {
	return composer.compose("DisciplinesList.md", messageData)
}

func (composer *MessageComposer) ComposeDisciplineScoresMessage(messageData models.DisciplinesScoresMessageData) (error, string) {
	return composer.compose("DisciplineScores.md", messageData)
}

func (composer *MessageComposer) ComposeScoreChanged(messageData models.ScoreChangedMessageData) (error, string) {
	if messageData.Previous.IsDeleted() {
		return composer.compose("ScoreChanged_created.md", messageData)
	} else if messageData.IsDeleted() {
		return composer.compose("ScoreChanged_deleted.md", messageData)
	} else {
		return composer.compose("ScoreChanged_changed.md", messageData)
	}
}

func (composer *MessageComposer) ComposeLogoutFinishedMessage() (error, string) {
	return composer.compose("LogoutFinished.md", nil)
}

func (composer *MessageComposer) compose(name string, data any) (error, string) {
	output := bytes.Buffer{}
	err := composer.templates.ExecuteTemplate(&output, name, data)
	return err, composer.postFilter(output.String())
}

func (composer *MessageComposer) getFunctionMap() template.FuncMap {
	return template.FuncMap{
		"renderScore": composer.renderScore,

		"incr": func(a int) int {
			return a + 1
		},

		"date": func(date time.Time) string {
			return date.Format("02.01.2006")
		},
	}
}

func (composer *MessageComposer) renderScore(score scoreApi.Score) string {
	if score.FirstScore != nil && score.SecondScore != nil {
		return composer.formatScore(score.FirstScore) + " та " + composer.formatScore(score.SecondScore)

	} else if score.FirstScore != nil {
		return composer.formatScore(score.FirstScore)

	} else if score.SecondScore != nil {
		return composer.formatScore(score.SecondScore)

	} else if score.IsAbsent {
		return "пропуск"

	} else {
		return "0"
	}
}

func (composer *MessageComposer) formatScore(score *float32) string {
	return strconv.FormatFloat(float64(*score), 'f', -1, 32)
}
