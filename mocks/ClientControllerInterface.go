// Code generated by mockery v2.28.1. DO NOT EDIT.

package mocks

import (
	context "context"

	events "github.com/kneu-messenger-pigeon/events"

	mock "github.com/stretchr/testify/mock"

	scoreApi "github.com/kneu-messenger-pigeon/score-api"

	sync "sync"
)

// ClientControllerInterface is an autogenerated mock type for the ClientControllerInterface type
type ClientControllerInterface struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, wg
func (_m *ClientControllerInterface) Execute(ctx context.Context, wg *sync.WaitGroup) {
	_m.Called(ctx, wg)
}

// LogoutFinishedAction provides a mock function with given fields: event
func (_m *ClientControllerInterface) LogoutFinishedAction(event *events.UserAuthorizedEvent) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(*events.UserAuthorizedEvent) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ScoreChangedAction provides a mock function with given fields: chatId, previousMessageId, disciplineScore, previousScore
func (_m *ClientControllerInterface) ScoreChangedAction(chatId string, previousMessageId string, disciplineScore *scoreApi.DisciplineScore, previousScore *scoreApi.Score) (error, string) {
	ret := _m.Called(chatId, previousMessageId, disciplineScore, previousScore)

	var r0 error
	var r1 string
	if rf, ok := ret.Get(0).(func(string, string, *scoreApi.DisciplineScore, *scoreApi.Score) (error, string)); ok {
		return rf(chatId, previousMessageId, disciplineScore, previousScore)
	}
	if rf, ok := ret.Get(0).(func(string, string, *scoreApi.DisciplineScore, *scoreApi.Score) error); ok {
		r0 = rf(chatId, previousMessageId, disciplineScore, previousScore)
	} else {
		r0 = ret.Error(0)
	}

	if rf, ok := ret.Get(1).(func(string, string, *scoreApi.DisciplineScore, *scoreApi.Score) string); ok {
		r1 = rf(chatId, previousMessageId, disciplineScore, previousScore)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// WelcomeAuthorizedAction provides a mock function with given fields: event
func (_m *ClientControllerInterface) WelcomeAuthorizedAction(event *events.UserAuthorizedEvent) error {
	ret := _m.Called(event)

	var r0 error
	if rf, ok := ret.Get(0).(func(*events.UserAuthorizedEvent) error); ok {
		r0 = rf(event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewClientControllerInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewClientControllerInterface creates a new instance of ClientControllerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewClientControllerInterface(t mockConstructorTestingTNewClientControllerInterface) *ClientControllerInterface {
	mock := &ClientControllerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
