// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	models "github.com/kneu-messenger-pigeon/client-framework/models"
	mock "github.com/stretchr/testify/mock"
)

// MessageComposerInterface is an autogenerated mock type for the MessageComposerInterface type
type MessageComposerInterface struct {
	mock.Mock
}

// ComposeDisciplineScoresMessage provides a mock function with given fields: messageData
func (_m *MessageComposerInterface) ComposeDisciplineScoresMessage(messageData models.DisciplinesScoresMessageData) (error, string) {
	ret := _m.Called(messageData)

	var r0 error
	var r1 string
	if rf, ok := ret.Get(0).(func(models.DisciplinesScoresMessageData) (error, string)); ok {
		return rf(messageData)
	}
	if rf, ok := ret.Get(0).(func(models.DisciplinesScoresMessageData) error); ok {
		r0 = rf(messageData)
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func(models.DisciplinesScoresMessageData) string); ok {
		r1 = rf(messageData)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// ComposeDisciplinesListMessage provides a mock function with given fields: messageData
func (_m *MessageComposerInterface) ComposeDisciplinesListMessage(messageData models.DisciplinesListMessageData) (error, string) {
	ret := _m.Called(messageData)

	var r0 error
	var r1 string
	if rf, ok := ret.Get(0).(func(models.DisciplinesListMessageData) (error, string)); ok {
		return rf(messageData)
	}
	if rf, ok := ret.Get(0).(func(models.DisciplinesListMessageData) error); ok {
		r0 = rf(messageData)
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func(models.DisciplinesListMessageData) string); ok {
		r1 = rf(messageData)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// ComposeLogoutFinishedMessage provides a mock function with given fields:
func (_m *MessageComposerInterface) ComposeLogoutFinishedMessage() (error, string) {
	ret := _m.Called()

	var r0 error
	var r1 string
	if rf, ok := ret.Get(0).(func() (error, string)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func() string); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// ComposeScoreChanged provides a mock function with given fields:
func (_m *MessageComposerInterface) ComposeScoreChanged() (error, string) {
	ret := _m.Called()

	var r0 error
	var r1 string
	if rf, ok := ret.Get(0).(func() (error, string)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func() string); ok {
		r1 = rf()
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// ComposeWelcomeAnonymousMessage provides a mock function with given fields: authUrl
func (_m *MessageComposerInterface) ComposeWelcomeAnonymousMessage(authUrl string) (error, string) {
	ret := _m.Called(authUrl)

	var r0 error
	var r1 string
	if rf, ok := ret.Get(0).(func(string) (error, string)); ok {
		return rf(authUrl)
	}
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(authUrl)
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func(string) string); ok {
		r1 = rf(authUrl)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// ComposeWelcomeAuthorizedMessage provides a mock function with given fields: messageData
func (_m *MessageComposerInterface) ComposeWelcomeAuthorizedMessage(messageData models.UserAuthorizedMessageData) (error, string) {
	ret := _m.Called(messageData)

	var r0 error
	var r1 string
	if rf, ok := ret.Get(0).(func(models.UserAuthorizedMessageData) (error, string)); ok {
		return rf(messageData)
	}
	if rf, ok := ret.Get(0).(func(models.UserAuthorizedMessageData) error); ok {
		r0 = rf(messageData)
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func(models.UserAuthorizedMessageData) string); ok {
		r1 = rf(messageData)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

type mockConstructorTestingTNewMessageComposerInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewMessageComposerInterface creates a new instance of MessageComposerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMessageComposerInterface(t mockConstructorTestingTNewMessageComposerInterface) *MessageComposerInterface {
	mock := &MessageComposerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
