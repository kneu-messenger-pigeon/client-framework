// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sync "sync"
)

// ExecutableInterface is an autogenerated mock type for the ExecutableInterface type
type ExecutableInterface struct {
	mock.Mock
}

// Execute provides a mock function with given fields: ctx, wg
func (_m *ExecutableInterface) Execute(ctx context.Context, wg *sync.WaitGroup) {
	_m.Called(ctx, wg)
}

type mockConstructorTestingTNewExecutableInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewExecutableInterface creates a new instance of ExecutableInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewExecutableInterface(t mockConstructorTestingTNewExecutableInterface) *ExecutableInterface {
	mock := &ExecutableInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
