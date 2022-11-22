package mocks

import mock "github.com/stretchr/testify/mock"

type AMQPPublisher struct {
	mock.Mock
}

// Create provides a mock function with given fields: value
func (_m *AMQPPublisher) Create(value string) {
	_m.Called(value)
}
