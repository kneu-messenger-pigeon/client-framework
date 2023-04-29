package framework

import (
	"bytes"
	"errors"
	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestUserRepository_SaveUser(t *testing.T) {
	t.Run("save_new", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		clientUserId := "test-id"
		student := Student{
			Name:       "Коваль Валера Павлович",
			Id:         99,
			LastName:   "Коваль",
			FirstName:  "Валера",
			MiddleName: "Павлович",
			Gender:     Student_MALE,
		}

		studentSerialized, _ := proto.Marshal(&student)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		redisMock.ExpectGetEx(clientUserId, UserExpiration).RedisNil()

		redisMock.ExpectTxPipeline()
		redisMock.ExpectSet(clientUserId, studentSerialized, UserExpiration).SetVal("OK")
		redisMock.ExpectSAdd(student.GetIdString(), clientUserId).SetVal(1)
		redisMock.ExpectExpire(student.GetIdString(), UserExpiration).SetVal(true)

		redisMock.ExpectTxPipelineExec()

		err := userRepository.SaveUser(clientUserId, &student)
		assert.NoError(t, err)
	})

	t.Run("replace_old", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		clientUserId := "test-id"
		previousStudent := Student{
			Name:       "Ткаченко Марія Андріївна",
			Id:         190,
			LastName:   "Ткаченко",
			FirstName:  "Марія",
			MiddleName: "Андріївна",
			Gender:     Student_FEMALE,
		}

		student := Student{
			Name:       "Коваль Валера Павлович",
			Id:         99,
			LastName:   "Коваль",
			FirstName:  "Валера",
			MiddleName: "Павлович",
			Gender:     Student_MALE,
		}

		previousStudentSerialized, _ := proto.Marshal(&previousStudent)
		studentSerialized, _ := proto.Marshal(&student)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		redisMock.ExpectGetEx(clientUserId, UserExpiration).SetVal(string(previousStudentSerialized))

		redisMock.ExpectTxPipeline()
		redisMock.ExpectDel(clientUserId).SetVal(1)
		redisMock.ExpectSRem(student.GetIdString(), clientUserId).SetVal(1)

		redisMock.ExpectSet(clientUserId, studentSerialized, UserExpiration).SetVal("OK")
		redisMock.ExpectSAdd(student.GetIdString(), clientUserId).SetVal(1)
		redisMock.ExpectExpire(student.GetIdString(), UserExpiration).SetVal(true)

		redisMock.ExpectTxPipelineExec()

		err := userRepository.SaveUser(clientUserId, &student)
		assert.NoError(t, err)
	})

	t.Run("delete_old", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		clientUserId := "test-id"
		previousStudent := Student{
			Name:       "Ткаченко Марія Андріївна",
			Id:         190,
			LastName:   "Ткаченко",
			FirstName:  "Марія",
			MiddleName: "Андріївна",
			Gender:     Student_FEMALE,
		}

		student := Student{
			Name:       "",
			Id:         0,
			LastName:   "",
			FirstName:  "",
			MiddleName: "",
			Gender:     Student_UNKNOWN,
		}

		previousStudentSerialized, _ := proto.Marshal(&previousStudent)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		redisMock.ExpectGetEx(clientUserId, UserExpiration).SetVal(string(previousStudentSerialized))

		redisMock.ExpectTxPipeline()
		redisMock.ExpectDel(clientUserId).SetVal(1)
		redisMock.ExpectSRem(student.GetIdString(), clientUserId).SetVal(1)

		redisMock.ExpectTxPipelineExec()

		err := userRepository.SaveUser(clientUserId, &student)
		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		expectedErr := errors.New("expected-test-error")

		clientUserId := "test-id"
		student := Student{
			Name:       "Коваль Валера Павлович",
			Id:         99,
			LastName:   "Коваль",
			FirstName:  "Валера",
			MiddleName: "Павлович",
			Gender:     Student_MALE,
		}

		studentSerialized, _ := proto.Marshal(&student)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		redisMock.ExpectGetEx(clientUserId, UserExpiration).RedisNil()

		redisMock.ExpectTxPipeline()
		redisMock.ExpectSet(clientUserId, studentSerialized, UserExpiration).SetVal("OK")
		redisMock.ExpectSAdd(student.GetIdString(), clientUserId).SetVal(1)
		redisMock.ExpectExpire(student.GetIdString(), UserExpiration).SetVal(true)

		redisMock.ExpectTxPipelineExec().SetErr(expectedErr)

		err := userRepository.SaveUser(clientUserId, &student)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestUserRepository_GetStudentUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		clientUserId := "test-id"

		student := &Student{
			Name:       "Коваль Валера Павлович",
			Id:         99,
			LastName:   "Коваль",
			FirstName:  "Валера",
			MiddleName: "Павлович",
			Gender:     Student_MALE,
		}

		studentSerialized, _ := proto.Marshal(student)

		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		redisMock.ExpectGetEx(clientUserId, UserExpiration).SetVal(string(studentSerialized))

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
		assertStudent(t, &Student{}, actualStudent)
	})
}

func assertStudent(t *testing.T, expected *Student, actual *Student) {
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
	t.Run("success", func(t *testing.T) {

		studentId := uint(100)

		redisClient, redisMock := redismock.NewClientMock()
		redisMock.MatchExpectationsInOrder(true)

		expectedIds := []string{
			"test-id-1",
			"test-id-2",
		}

		redisMock.ExpectSMembers("100").SetVal(expectedIds)

		userRepository := UserRepository{
			out:   &bytes.Buffer{},
			redis: redisClient,
		}

		actualIds := userRepository.GetClientUserIds(studentId)
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
