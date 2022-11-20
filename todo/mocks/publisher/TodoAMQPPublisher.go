package mocks

import (
	mock "github.com/stretchr/testify/mock"
)

type TodoAMQPPublisher struct {
	mock.Mock
}

func (_m *TodoAMQPPublisher) Create(value string) {}
