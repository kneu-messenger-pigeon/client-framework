package models

import (
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
	"time"
)

type WelcomeAnonymousMessageData struct {
	AuthUrl  string
	ExpireAt time.Time
}

type StudentMessageData struct {
	NamePrefix string
	Name       string
}

type UserAuthorizedMessageData struct {
	StudentMessageData
}

type DisciplinesListMessageData struct {
	StudentMessageData
	Disciplines scoreApi.DisciplineScoreResults
}

type DisciplinesScoresMessageData struct {
	StudentMessageData
	Discipline scoreApi.DisciplineScoreResult
}

type ScoreChangedMessageData struct {
	scoreApi.Discipline
	scoreApi.Score
	Previous scoreApi.Score
}

var studentNamePrefixMap = map[Student_GenderType]string{
	Student_FEMALE:  "Пані",
	Student_MALE:    "Пане",
	Student_UNKNOWN: "Пане",
}

func NewStudentMessageData(student *Student) StudentMessageData {
	return StudentMessageData{
		NamePrefix: studentNamePrefixMap[student.Gender],
		Name:       student.FirstName,
	}
}
