package framework

import (
	"context"
	"github.com/kneu-messenger-pigeon/client-framework/models"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/proto"
	"io"
	"strconv"
	"time"
)

type UserRepository struct {
	out   io.Writer
	redis redis.UniversalClient
}

const UserExpiration = time.Hour * 24 * 30 * 7 // 7 months, 210 days

const RedisBackgroundSaveInProgress = "ERR Background save already in progress"

func (repository *UserRepository) SaveUser(clientUserId string, student *models.Student) (err error) {
	previousStudent := repository.GetStudent(clientUserId)

	studentSerialized, err := proto.Marshal(student)
	ctx := context.Background()
	if err == nil {
		clientUserKey := repository.getClientUserKey(clientUserId)
		pipe := repository.redis.TxPipeline()

		if previousStudent.Id != student.Id && previousStudent.Id != 0 {
			pipe.SRem(ctx, repository.getStudentKey(previousStudent.Id), clientUserId)
		}

		if student.Id == 0 {
			pipe.Del(ctx, clientUserKey)
		} else {
			pipe.Set(ctx, clientUserKey, studentSerialized, UserExpiration)
			newStudentKey := repository.getStudentKey(student.Id)
			pipe.SAdd(ctx, newStudentKey, clientUserId)
			pipe.Expire(ctx, newStudentKey, UserExpiration)
		}

		_, err = pipe.Exec(ctx)
	}

	return err
}

func (repository *UserRepository) Commit() error {
	err := repository.redis.BgSave(context.Background()).Err()
	if err != nil && err.Error() == RedisBackgroundSaveInProgress {
		err = nil
	}

	return err
}

func (repository *UserRepository) GetStudent(clientUserId string) *models.Student {
	ctx := context.Background()
	studentSerialized, _ := repository.redis.GetEx(
		ctx, repository.getClientUserKey(clientUserId),
		UserExpiration,
	).Bytes()

	student := &models.Student{}
	if studentSerialized != nil && len(studentSerialized) > 0 {
		_ = proto.Unmarshal(studentSerialized, student)
	}

	if student.Id != 0 {
		repository.redis.Expire(ctx, repository.getStudentKey(student.Id), UserExpiration)
	}

	return student
}

func (repository *UserRepository) GetClientUserIds(studentId uint) []string {
	if studentId != 0 {
		result := repository.redis.SMembers(
			context.Background(),
			repository.getStudentKey(uint32(studentId)),
		)

		if result.Err() == nil {
			return result.Val()
		}
	}

	return []string{}
}

func (repository *UserRepository) getStudentKey(studentId uint32) string {
	return "st" + strconv.FormatUint(uint64(studentId), 10)
}

func (repository *UserRepository) getClientUserKey(clientUserId string) string {
	return "cu" + clientUserId
}
