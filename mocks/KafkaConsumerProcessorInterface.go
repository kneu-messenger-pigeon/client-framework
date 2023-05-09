// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	sync "sync"
)

// KafkaConsumerProcessorInterface is an autogenerated mock type for the KafkaConsumerProcessorInterface type
type KafkaConsumerProcessorInterface struct {
	mock.Mock
}

// Disable provides a mock function with given fields:
func (_m *KafkaConsumerProcessorInterface) Disable() {
	_m.Called()
}

// Execute provides a mock function with given fields: ctx, wg
func (_m *KafkaConsumerProcessorInterface) Execute(ctx context.Context, wg *sync.WaitGroup) {
	_m.Called(ctx, wg)
}

type mockConstructorTestingTNewKafkaConsumerProcessorInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewKafkaConsumerProcessorInterface creates a new instance of KafkaConsumerProcessorInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewKafkaConsumerProcessorInterface(t mockConstructorTestingTNewKafkaConsumerProcessorInterface) *KafkaConsumerProcessorInterface {
	mock := &KafkaConsumerProcessorInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
