package framework

import (
	"context"
	"github.com/kneu-messenger-pigeon/client-framework/models"
)

type UserRepositoryInterface interface {
	SaveUser(clientUserId string, student *models.Student) (err error, hasChanges bool)
	GetStudent(clientUserId string) *models.Student
	GetClientUserIds(studentId uint) []string
	Commit() error
	GetUserCount(ctx context.Context) (redisUserCount uint64, err error)
}
