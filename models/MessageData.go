package models

import (
	scoreApi "github.com/kneu-messenger-pigeon/score-api"
)

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

func NewStudentMessageData(student *Student) StudentMessageData {
	return StudentMessageData{
		NamePrefix: student.GetNamePrefix(),
		Name:       student.FirstName,
	}
}
