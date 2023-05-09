package framework

import "github.com/kneu-messenger-pigeon/client-framework/models"

type UserRepositoryInterface interface {
	SaveUser(clientUserId string, student *models.Student) error
	GetStudent(clientUserId string) *models.Student
	GetClientUserIds(studentId uint) []string
	Commit() error
}
