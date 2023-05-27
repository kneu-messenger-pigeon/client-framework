package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewStudentMessageData(t *testing.T) {
	t.Run("male", func(t *testing.T) {
		student := &Student{
			Id:         uint32(999),
			LastName:   "Потапенко",
			FirstName:  "Андрій",
			MiddleName: "Петрович",
			Gender:     Student_MALE,
		}

		studentMessageData := NewStudentMessageData(student)

		assert.NotEmpty(t, studentMessageData)
		assert.Equal(t, "Пане", studentMessageData.NamePrefix)
		assert.Equal(t, "Андрій", studentMessageData.Name)
	})

	t.Run("female", func(t *testing.T) {
		student := &Student{
			Id:         uint32(999),
			LastName:   "Потапенко",
			FirstName:  "Марія",
			MiddleName: "Петрівна",
			Gender:     Student_FEMALE,
		}

		studentMessageData := NewStudentMessageData(student)

		assert.NotEmpty(t, studentMessageData)
		assert.Equal(t, "Пані", studentMessageData.NamePrefix)
		assert.Equal(t, "Марія", studentMessageData.Name)
	})

}
