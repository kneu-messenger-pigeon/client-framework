package framework

import (
	"bytes"
	"context"
	"errors"
	"github.com/alicebob/miniredis/v2"
	"github.com/go-redis/redismock/v9"
	"github.com/kneu-messenger-pigeon/client-framework/models"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestUserRepository_SaveUser(t *testing.T) {
	t.Run("save_new", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		clientUserId := "test-id"
		student := models.Student{
			Name:       "Коваль Валера Павлович",
			Id:         99,
			LastName:   "Коваль",
			FirstName:  "Валера",
			MiddleName: "Павлович",
			Gender:     models.Student_MALE,
		}

		studentSerialized, _ := proto.Marshal(&student)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		studentKey := userRepository.getStudentKey(student.Id)
		clientUserKey := userRepository.getClientUserKey(clientUserId)
		redisMock.ExpectGetEx(clientUserKey, UserExpiration).RedisNil()

		redisMock.ExpectTxPipeline()
		redisMock.ExpectSet(clientUserKey, studentSerialized, UserExpiration).SetVal("OK")
		redisMock.ExpectSAdd(studentKey, clientUserId).SetVal(1)
		redisMock.ExpectExpire(studentKey, UserExpiration).SetVal(true)

		redisMock.ExpectTxPipelineExec()

		err, hasChanges := userRepository.SaveUser(clientUserId, &student)
		assert.True(t, hasChanges)
		assert.NoError(t, err)
		assert.NoError(t, redisMock.ExpectationsWereMet())
	})

	t.Run("replace_old", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		clientUserId := "test-id"
		previousStudent := models.Student{
			Name:       "Ткаченко Марія Андріївна",
			Id:         190,
			LastName:   "Ткаченко",
			FirstName:  "Марія",
			MiddleName: "Андріївна",
			Gender:     models.Student_FEMALE,
		}

		student := models.Student{
			Name:       "Коваль Валера Павлович",
			Id:         99,
			LastName:   "Коваль",
			FirstName:  "Валера",
			MiddleName: "Павлович",
			Gender:     models.Student_MALE,
		}

		previousStudentSerialized, _ := proto.Marshal(&previousStudent)
		studentSerialized, _ := proto.Marshal(&student)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}
		previousStudentKey := userRepository.getStudentKey(previousStudent.Id)
		newStudentKey := userRepository.getStudentKey(student.Id)
		clientUserKey := userRepository.getClientUserKey(clientUserId)

		redisMock.ExpectGetEx(clientUserKey, UserExpiration).SetVal(string(previousStudentSerialized))

		redisMock.ExpectTxPipeline()
		redisMock.ExpectSRem(previousStudentKey, clientUserId).SetVal(1)

		redisMock.ExpectSet(clientUserKey, studentSerialized, UserExpiration).SetVal("OK")
		redisMock.ExpectSAdd(newStudentKey, clientUserId).SetVal(1)
		redisMock.ExpectExpire(newStudentKey, UserExpiration).SetVal(true)

		redisMock.ExpectTxPipelineExec()

		err, hasChanges := userRepository.SaveUser(clientUserId, &student)
		assert.True(t, hasChanges)
		assert.NoError(t, err)
	})

	t.Run("delete_old", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		clientUserId := "test-id"
		previousStudent := models.Student{
			Name:       "Ткаченко Марія Андріївна",
			Id:         190,
			LastName:   "Ткаченко",
			FirstName:  "Марія",
			MiddleName: "Андріївна",
			Gender:     models.Student_FEMALE,
		}

		student := models.Student{
			Name:       "",
			Id:         0,
			LastName:   "",
			FirstName:  "",
			MiddleName: "",
			Gender:     models.Student_UNKNOWN,
		}

		previousStudentSerialized, _ := proto.Marshal(&previousStudent)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}
		previousStudentKey := userRepository.getStudentKey(previousStudent.Id)
		clientUserKey := userRepository.getClientUserKey(clientUserId)

		redisMock.ExpectGetEx(clientUserKey, UserExpiration).SetVal(string(previousStudentSerialized))

		redisMock.ExpectTxPipeline()
		redisMock.ExpectSRem(previousStudentKey, clientUserId).SetVal(1)
		redisMock.ExpectDel(clientUserKey).SetVal(1)

		redisMock.ExpectTxPipelineExec()

		err, hasChanges := userRepository.SaveUser(clientUserId, &student)
		assert.True(t, hasChanges)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		expectedErr := errors.New("expected-test-error")

		clientUserId := "test-id"
		student := models.Student{
			Name:       "Коваль Валера Павлович",
			Id:         99,
			LastName:   "Коваль",
			FirstName:  "Валера",
			MiddleName: "Павлович",
			Gender:     models.Student_MALE,
		}

		studentSerialized, _ := proto.Marshal(&student)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}
		studentKey := userRepository.getStudentKey(student.Id)
		clientUserKey := userRepository.getClientUserKey(clientUserId)

		redisMock.ExpectGetEx(clientUserKey, UserExpiration).RedisNil()

		redisMock.ExpectTxPipeline()
		redisMock.ExpectSet(clientUserKey, studentSerialized, UserExpiration).SetVal("OK")
		redisMock.ExpectSAdd(studentKey, clientUserId).SetVal(1)
		redisMock.ExpectExpire(studentKey, UserExpiration).SetVal(true)

		redisMock.ExpectTxPipelineExec().SetErr(expectedErr)

		err, hasChanges := userRepository.SaveUser(clientUserId, &student)
		assert.True(t, hasChanges)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("saveNewAndRetrieveAndReplaceWithOtherStudentAndDelete", func(t *testing.T) {
		expectedClientUserId := "test-id"
		expectedStudent1 := &models.Student{
			Name:       "Коваль Валера Павлович",
			Id:         99,
			LastName:   "Коваль",
			FirstName:  "Валера",
			MiddleName: "Павлович",
			Gender:     models.Student_MALE,
		}

		expectedStudent2 := &models.Student{
			Name:       "Ткаченко Юлія Андрійвна",
			Id:         235,
			LastName:   "Ткаченко",
			FirstName:  "Юлія",
			MiddleName: "Андрійвна",
			Gender:     models.Student_FEMALE,
		}

		emptyStudent := &models.Student{}

		userRepository := UserRepository{
			out: &bytes.Buffer{},
			redis: redis.NewClient(&redis.Options{
				Network: "tcp",
				Addr:    miniredis.RunT(t).Addr(),
			}),
		}

		// save student for new client user
		err, hasChanges := userRepository.SaveUser(expectedClientUserId, expectedStudent1)
		assert.NoError(t, err)
		assert.True(t, hasChanges)

		err, hasChanges = userRepository.SaveUser(expectedClientUserId, expectedStudent1)
		assert.False(t, hasChanges)

		actualStudent := userRepository.GetStudent(expectedClientUserId)
		assert.Equal(t, expectedStudent1.String(), actualStudent.String())

		actualClientIds := userRepository.GetClientUserIds(uint(expectedStudent1.Id))
		assert.Len(t, actualClientIds, 1)
		assert.Equal(t, expectedClientUserId, actualClientIds[0])

		// replace student
		err, hasChanges = userRepository.SaveUser(expectedClientUserId, expectedStudent2)
		assert.True(t, hasChanges)

		actualStudent = userRepository.GetStudent(expectedClientUserId)
		assert.Equal(t, expectedStudent2.String(), actualStudent.String())

		actualClientIds = userRepository.GetClientUserIds(uint(expectedStudent1.Id))
		assert.Len(t, actualClientIds, 0)

		actualClientIds = userRepository.GetClientUserIds(uint(expectedStudent2.Id))
		assert.Len(t, actualClientIds, 1)
		assert.Equal(t, expectedClientUserId, actualClientIds[0])

		// delete student
		err, hasChanges = userRepository.SaveUser(expectedClientUserId, emptyStudent)
		assert.True(t, hasChanges)
		assert.NoError(t, err)

		actualStudent = userRepository.GetStudent(expectedClientUserId)
		assert.Nil(t, actualStudent)

		actualClientIds = userRepository.GetClientUserIds(uint(expectedStudent1.Id))
		assert.Len(t, actualClientIds, 0)
		actualClientIds = userRepository.GetClientUserIds(uint(expectedStudent2.Id))
		assert.Len(t, actualClientIds, 0)
	})
}

func TestUserRepository_GetStudent(t *testing.T) {
	t.Run("clientKey", func(t *testing.T) {
		userRepository := UserRepository{}
		clientKey := userRepository.getClientUserKey("1u1")
		assert.Equal(t, "cu1u1", clientKey)
	})

	t.Run("success", func(t *testing.T) {
		clientUserId := "test-id"

		student := &models.Student{
			Name:       "Коваль Валера Павлович",
			Id:         99,
			LastName:   "Коваль",
			FirstName:  "Валера",
			MiddleName: "Павлович",
			Gender:     models.Student_MALE,
		}

		studentSerialized, _ := proto.Marshal(student)

		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}
		clientUserKey := userRepository.getClientUserKey(clientUserId)

		redisMock.ExpectGetEx(clientUserKey, UserExpiration).SetVal(string(studentSerialized))

		actualStudent := userRepository.GetStudent(clientUserId)
		assertStudent(t, student, actualStudent)
	})

	t.Run("not_exists", func(t *testing.T) {
		clientUserId := "test-id"

		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		redisMock.ExpectGetEx(clientUserId, UserExpiration).RedisNil()

		actualStudent := userRepository.GetStudent(clientUserId)
		assert.Nil(t, actualStudent)
	})
}

func assertStudent(t *testing.T, expected *models.Student, actual *models.Student) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.LastName, actual.LastName)
	assert.Equal(t, expected.FirstName, actual.FirstName)
	assert.Equal(t, expected.MiddleName, actual.MiddleName)
	assert.Equal(t, expected.Gender, actual.Gender)
}

func TestUserRepository_Commit(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		redisMock.ExpectBgSave().SetVal("OK")

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		err := userRepository.Commit()
		assert.NoError(t, err)
	})

	t.Run("save_in_progress", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		redisMock.ExpectBgSave().SetErr(errors.New(RedisBackgroundSaveInProgress))

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		err := userRepository.Commit()
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		expectedError := errors.New("Test expected error")

		redisMock.ExpectBgSave().SetErr(expectedError)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		err := userRepository.Commit()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestUserRepository_GetClientUserIds(t *testing.T) {
	t.Run("getStudentKey", func(t *testing.T) {
		userRepository := UserRepository{}
		studentKey := userRepository.getStudentKey(123)
		assert.Equal(t, "st123", studentKey)
	})

	t.Run("success", func(t *testing.T) {
		studentId := uint(100)

		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		expectedIds := []string{
			"test-id-1",
			"test-id-2",
		}

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}
		redisMock.ExpectSMembers(
			userRepository.getStudentKey(uint32(studentId)),
		).SetVal(expectedIds)

		actualIds := userRepository.GetClientUserIds(studentId)
		assert.NoError(t, redisMock.ExpectationsWereMet())
		assert.Equal(t, expectedIds, actualIds)
	})

	t.Run("empty", func(t *testing.T) {
		studentId := uint(100)

		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		redisMock.ExpectSMembers("100").RedisNil()

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		actualIds := userRepository.GetClientUserIds(studentId)
		assert.Equal(t, []string{}, actualIds)
	})
}

func TestUserRepository_GetUserCount(t *testing.T) {
	t.Run("succes", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		match := ClientUserPrefix + "*"

		returnKeysBatch1 := []string{
			"cu1",
			"cu2",
			"cu3",
			"cu4",
		}

		returnKeysBatch2 := []string{
			"cu9",
			"cu10",
			"cu11",
			"cu12",
			"cu14",
			"cu16",
		}

		nextCursor := uint64(110)

		redisMock.ExpectScan(0, match, UserScanBatchSize).SetVal(returnKeysBatch1, nextCursor)
		redisMock.ExpectScan(nextCursor, match, UserScanBatchSize).SetVal(returnKeysBatch2, 0)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		count, err := userRepository.GetUserCount(context.Background())

		expectCount := uint64(len(returnKeysBatch1) + len(returnKeysBatch2))

		assert.NoError(t, err)
		assert.Equal(t, expectCount, count)
	})
}
