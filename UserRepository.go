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

const ClientUserPrefix = "cu"

const UserScanBatchSize = 500

func (repository *UserRepository) SaveUser(clientUserId string, student *models.Student) (error, bool) {
	previousStudent := repository.GetStudent(clientUserId)

	studentSerialized, err := proto.Marshal(student)
	ctx := context.Background()

	hasChanges := previousStudent != student
	if previousStudent != nil && student != nil {
		hasChanges = previousStudent.Id != student.Id
	}

	if err == nil {
		clientUserKey := repository.getClientUserKey(clientUserId)
		pipe := repository.redis.TxPipeline()

		if hasChanges && previousStudent != nil {
			pipe.SRem(ctx, repository.getStudentKey(previousStudent.Id), clientUserId)
		}

		if student == nil || student.Id == 0 {
			pipe.Del(ctx, clientUserKey)
		} else {
			pipe.Set(ctx, clientUserKey, studentSerialized, UserExpiration)
			newStudentKey := repository.getStudentKey(student.Id)
			pipe.SAdd(ctx, newStudentKey, clientUserId)
			pipe.Expire(ctx, newStudentKey, UserExpiration)
		}

		_, err = pipe.Exec(ctx)
	}

	return err, hasChanges
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

	if studentSerialized != nil && len(studentSerialized) > 0 {
		student := models.Student{}
		_ = proto.Unmarshal(studentSerialized, &student)
		if student.Id != 0 {
			repository.redis.Expire(ctx, repository.getStudentKey(student.Id), UserExpiration)
			return &student
		}
	}

	return nil
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
	return ClientUserPrefix + clientUserId
}

func (repository *UserRepository) GetUserCount(ctx context.Context) (redisUserCount uint64, err error) {
	match := ClientUserPrefix + "*"

	var cursor uint64

	var keys []string
	for err == nil {
		keys, cursor, err = repository.redis.Scan(ctx, cursor, match, UserScanBatchSize).Result()
		redisUserCount += uint64(len(keys))
		if cursor == 0 {
			break
		}
	}

	return
}
