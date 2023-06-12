package framework

import (
	"github.com/kneu-messenger-pigeon/client-framework/models"
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
	"time"
)

var composeTestStudentData = models.StudentMessageData{
	NamePrefix: "Пане",
	Name:       "Микита",
}

func TestMessageComposer_ComposeWelcomeAnonymousMessage(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		var composer MessageComposerInterface
		composer = NewMessageComposer(MessageComposerConfig{})

		authUrl := "https://example.com/auth"

		err, message := composer.ComposeWelcomeAnonymousMessage(authUrl)

		assert.NoError(t, err)
		assert.Contains(t, message, authUrl)
	})
}

func TestMessageComposer_ComposeWelcomeAuthorizedMessage(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		composer := NewMessageComposer(MessageComposerConfig{})

		messageData := models.UserAuthorizedMessageData{
			StudentMessageData: composeTestStudentData,
		}

		err, message := composer.ComposeWelcomeAuthorizedMessage(messageData)
		assert.NoError(t, err)
		assert.Contains(t, message, messageData.NamePrefix)
		assert.Contains(t, message, messageData.Name)
	})
}

func TestMessageComposer_ComposeDisciplinesListMessage(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		disciplines := scoreApi.DisciplineScoreResults{
			{
				Discipline: scoreApi.Discipline{
					Id:   110,
					Name: "Наноекономіка",
				},
				ScoreRating: scoreApi.ScoreRating{
					Total:         30,
					MinTotal:      20,
					MaxTotal:      40,
					Rating:        7,
					StudentsCount: 12,
				},
				Scores: []scoreApi.Score{},
			},
			{
				Discipline: scoreApi.Discipline{
					Id:   220,
					Name: "Культура кави у бізнесі",
				},
				ScoreRating: scoreApi.ScoreRating{
					Total:         18,
					MinTotal:      16,
					MaxTotal:      19,
					Rating:        1,
					StudentsCount: 8,
				},
				Scores: []scoreApi.Score{},
			},
		}

		messageData := models.DisciplinesListMessageData{
			StudentMessageData: composeTestStudentData,
			Disciplines:        disciplines,
		}

		composer := NewMessageComposer(MessageComposerConfig{})
		err, message := composer.ComposeDisciplinesListMessage(messageData)
		assert.NoError(t, err)
		assert.Contains(t, message, messageData.NamePrefix)
		assert.Contains(t, message, messageData.Name)

		for _, discipline := range disciplines {
			assert.Contains(t, message, discipline.Discipline.Name)
			assert.Contains(t, message, _formatFloat(discipline.ScoreRating.Total))
			assert.Contains(t, message, strconv.Itoa(discipline.ScoreRating.Rating))
			assert.Contains(t, message, strconv.Itoa(discipline.ScoreRating.StudentsCount))
		}

		assert.Contains(t, message, "офіційному журналі успішності КНЕУ")
		assert.Contains(t, message, "https://cutt.ly/Dekanat")
	})
}

func TestMessageComposer_ComposeDisciplineScoresMessage(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		discipline := scoreApi.DisciplineScoreResult{
			Discipline: scoreApi.Discipline{
				Id:   110,
				Name: "Наноекономіка",
			},
			ScoreRating: scoreApi.ScoreRating{
				Total:         41,
				MinTotal:      8,
				MaxTotal:      48,
				Rating:        7,
				StudentsCount: 12,
			},
			Scores: []scoreApi.Score{
				{
					Lesson: scoreApi.Lesson{
						Id:   451,
						Date: time.Date(2023, 5, 15, 0, 0, 0, 0, time.Local),
						Type: scoreApi.LessonType{
							Id:        20,
							ShortName: "Лек",
							LongName:  "Лекція",
						},
					},
					FirstScore:  floatPointer(3),
					SecondScore: nil,
					IsAbsent:    false,
				},
				{
					Lesson: scoreApi.Lesson{
						Id:   453,
						Date: time.Date(2023, 5, 18, 0, 0, 0, 0, time.Local),
						Type: scoreApi.LessonType{
							Id:        25,
							ShortName: "Сем",
							LongName:  "Семінар",
						},
					},
					FirstScore:  nil,
					SecondScore: nil,
					IsAbsent:    true,
				},
				{
					Lesson: scoreApi.Lesson{
						Id:   453,
						Date: time.Date(2023, 5, 18, 0, 0, 0, 0, time.Local),
						Type: scoreApi.LessonType{
							Id:        25,
							ShortName: "Сем",
							LongName:  "Семінар",
						},
					},
					FirstScore:  nil,
					SecondScore: floatPointer(6),
					IsAbsent:    true,
				},
				{
					Lesson: scoreApi.Lesson{
						Id:   453,
						Date: time.Date(2023, 5, 19, 0, 0, 0, 0, time.Local),
						Type: scoreApi.LessonType{
							Id:        25,
							ShortName: "Лаб",
							LongName:  "Лабороторна робота",
						},
					},
					FirstScore:  nil,
					SecondScore: nil,
					IsAbsent:    false,
				},
				{
					Lesson: scoreApi.Lesson{
						Id:   456,
						Date: time.Date(2023, 5, 20, 0, 0, 0, 0, time.Local),
						Type: scoreApi.LessonType{
							Id:        28,
							ShortName: "Реф",
							LongName:  "Реферат",
						},
					},
					FirstScore:  floatPointer(2),
					SecondScore: floatPointer(1),
					IsAbsent:    false,
				},
			},
		}

		messageData := models.DisciplinesScoresMessageData{
			StudentMessageData: composeTestStudentData,
			Discipline:         discipline,
		}

		composer := NewMessageComposer(MessageComposerConfig{})
		err, message := composer.ComposeDisciplineScoresMessage(messageData)

		assert.NoError(t, err)
		assert.Contains(t, message, discipline.Discipline.Name)
		assert.Contains(t, message, _formatFloat(discipline.ScoreRating.MinTotal))
		assert.Contains(t, message, _formatFloat(discipline.ScoreRating.MaxTotal))

		assert.Equal(t, 1, strings.Count(message, "15.05.2023"))
		assert.Equal(t, 2, strings.Count(message, "18.05.2023"))
		assert.Equal(t, 1, strings.Count(message, "19.05.2023"))
		assert.Equal(t, 1, strings.Count(message, "20.05.2023"))

		assert.Equal(t, 1, strings.Count(message, "*3*"))
		assert.Equal(t, 1, strings.Count(message, "*пропуск*"))
		assert.Equal(t, 1, strings.Count(message, "*6*"))
		assert.Equal(t, 1, strings.Count(message, "*0*"))
		assert.Equal(t, 1, strings.Count(message, "*2 та 1*"))

		for _, score := range discipline.Scores {
			assert.Contains(t, message, score.Lesson.Type.LongName)
		}
	})
}

func TestMessageComposer_ComposeScoreChanged(t *testing.T) {
	discipline := scoreApi.Discipline{
		Id:   110,
		Name: "Наноекономіка",
	}

	t.Run("created_score", func(t *testing.T) {
		composer := NewMessageComposer(MessageComposerConfig{})

		messageData := models.ScoreChangedMessageData{
			Discipline: discipline,
			Score: scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   456,
					Date: time.Date(2023, 5, 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        28,
						ShortName: "Реф",
						LongName:  "Реферат",
					},
				},
				FirstScore:  floatPointer(5.5),
				SecondScore: floatPointer(1),
				IsAbsent:    false,
			},
			Previous: scoreApi.Score{},
		}

		err, message := composer.ComposeScoreChanged(messageData)

		assert.NoError(t, err)
		assert.Equal(t, "Новий запис: Наноекономіка, заняття 12.05.2023 _Реферат_: 5.5 та 1", message)
	})

	t.Run("deleted_score", func(t *testing.T) {
		composer := NewMessageComposer(MessageComposerConfig{})

		messageData := models.ScoreChangedMessageData{
			Discipline: discipline,
			Score: scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   456,
					Date: time.Date(2023, 5, 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        28,
						ShortName: "Реф",
						LongName:  "Реферат",
					},
				},
			},
			Previous: scoreApi.Score{
				FirstScore: floatPointer(1.5),
				IsAbsent:   false,
			},
		}
		err, message := composer.ComposeScoreChanged(messageData)

		assert.NoError(t, err)
		assert.Equal(t, "Видалено запис: Наноекономіка, заняття 12.05.2023 _Реферат_: ~1.5~", message)
	})

	t.Run("changed_score", func(t *testing.T) {
		composer := NewMessageComposer(MessageComposerConfig{})

		messageData := models.ScoreChangedMessageData{
			Discipline: discipline,
			Score: scoreApi.Score{
				Lesson: scoreApi.Lesson{
					Id:   456,
					Date: time.Date(2023, 5, 12, 0, 0, 0, 0, time.Local),
					Type: scoreApi.LessonType{
						Id:        28,
						ShortName: "Реф",
						LongName:  "Реферат",
					},
				},
				FirstScore:  floatPointer(5.5),
				SecondScore: floatPointer(1),
				IsAbsent:    false,
			},
			Previous: scoreApi.Score{
				FirstScore:  floatPointer(2.5),
				SecondScore: floatPointer(2.0),
			},
		}
		err, message := composer.ComposeScoreChanged(messageData)

		assert.NoError(t, err)
		assert.Equal(t, "Змінено запис: Наноекономіка, заняття 12.05.2023 _Реферат_: 5.5 та 1 (було ~2.5 та 2~)", message)
	})
}

func TestMessageComposer_ComposeLogoutFinishedMessage(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		composer := NewMessageComposer(MessageComposerConfig{})
		err, message := composer.ComposeLogoutFinishedMessage()

		assert.NoError(t, err)
		assert.Contains(t, message, "зупинено")
	})
}

func TestNewMessageComposer_PostFilter(t *testing.T) {
	replacedString := "REPLACED"

	composer := NewMessageComposer(MessageComposerConfig{})
	composer.SetPostFilter(func(i string) string {
		return replacedString
	})

	err, actualMessage := composer.ComposeLogoutFinishedMessage()

	assert.NoError(t, err)
	assert.Equal(t, replacedString, actualMessage)
}

func _formatFloat(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', 0, 32)
}

func floatPointer(value float32) *float32 {
	return &value
}
