package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStudent_GetIdString(t *testing.T) {
	student := &Student{
		Id:         uint32(999),
		LastName:   "Потапенко",
		FirstName:  "Андрій",
		MiddleName: "Петрович",
		Gender:     Student_MALE,
	}

	assert.Equal(t, "999", student.GetIdString())
}
