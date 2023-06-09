// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	models "github.com/kneu-messenger-pigeon/client-framework/models"
	mock "github.com/stretchr/testify/mock"
)

// ScoreChangedMessageIdStorageInterface is an autogenerated mock type for the ScoreChangedMessageIdStorageInterface type
type ScoreChangedMessageIdStorageInterface struct {
	mock.Mock
}

// GetAll provides a mock function with given fields: studentId, lessonId
func (_m *ScoreChangedMessageIdStorageInterface) GetAll(studentId uint, lessonId uint) models.ScoreChangedMessageMap {
	ret := _m.Called(studentId, lessonId)

	var r0 models.ScoreChangedMessageMap
	if rf, ok := ret.Get(0).(func(uint, uint) models.ScoreChangedMessageMap); ok {
		r0 = rf(studentId, lessonId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(models.ScoreChangedMessageMap)
		}
	}

	return r0
}

// Set provides a mock function with given fields: studentId, lessonId, chatId, messageId
func (_m *ScoreChangedMessageIdStorageInterface) Set(studentId uint, lessonId uint, chatId string, messageId string) {
	_m.Called(studentId, lessonId, chatId, messageId)
}

type mockConstructorTestingTNewScoreChangedMessageIdStorageInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewScoreChangedMessageIdStorageInterface creates a new instance of ScoreChangedMessageIdStorageInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewScoreChangedMessageIdStorageInterface(t mockConstructorTestingTNewScoreChangedMessageIdStorageInterface) *ScoreChangedMessageIdStorageInterface {
	mock := &ScoreChangedMessageIdStorageInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
