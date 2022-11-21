package mocks

import mock "github.com/stretchr/testify/mock"

type TodoAMQPPublisher struct {
	mock.Mock
}

// Create provides a mock function with given fields: value
func (_m *TodoAMQPPublisher) Create(value string) {
	_m.Called(value)
}
