// Code generated by mockery v2.10.4. DO NOT EDIT.

package mocks

import (
	context "context"
	models "go-rengan/todo/models"

	mock "github.com/stretchr/testify/mock"
)

// MongoTodoRepository is an autogenerated mock type for the MongoTodoRepository type
type MongoTodoRepository struct {
	mock.Mock
}

// CountFindAll provides a mock function with given fields: ctx, keyword
func (_m *MongoTodoRepository) CountFindAll(ctx context.Context, keyword string) (int, error) {
	ret := _m.Called(ctx, keyword)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, keyword)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, keyword)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CountFindByID provides a mock function with given fields: ctx, id
func (_m *MongoTodoRepository) CountFindByID(ctx context.Context, id string) (int, error) {
	ret := _m.Called(ctx, id)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context, string) int); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: ctx, id
func (_m *MongoTodoRepository) Delete(ctx context.Context, id string) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindAll provides a mock function with given fields: ctx, keyword, limit, offset
func (_m *MongoTodoRepository) FindAll(ctx context.Context, keyword string, limit int, offset int) ([]*models.Todo, error) {
	ret := _m.Called(ctx, keyword, limit, offset)

	var r0 []*models.Todo
	if rf, ok := ret.Get(0).(func(context.Context, string, int, int) []*models.Todo); ok {
		r0 = rf(ctx, keyword, limit, offset)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Todo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, int, int) error); ok {
		r1 = rf(ctx, keyword, limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FindById provides a mock function with given fields: ctx, id
func (_m *MongoTodoRepository) FindById(ctx context.Context, id string) (*models.Todo, error) {
	ret := _m.Called(ctx, id)

	var r0 *models.Todo
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Todo); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Todo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Store provides a mock function with given fields: ctx, value
func (_m *MongoTodoRepository) Store(ctx context.Context, value *models.Todo) (*models.Todo, error) {
	ret := _m.Called(ctx, value)

	var r0 *models.Todo
	if rf, ok := ret.Get(0).(func(context.Context, *models.Todo) *models.Todo); ok {
		r0 = rf(ctx, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Todo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *models.Todo) error); ok {
		r1 = rf(ctx, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: ctx, id, value
func (_m *MongoTodoRepository) Update(ctx context.Context, id string, value *models.Todo) (*models.Todo, error) {
	ret := _m.Called(ctx, id, value)

	var r0 *models.Todo
	if rf, ok := ret.Get(0).(func(context.Context, string, *models.Todo) *models.Todo); ok {
		r0 = rf(ctx, id, value)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Todo)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, *models.Todo) error); ok {
		r1 = rf(ctx, id, value)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
