package service_test

import (
	"context"
	tracing "go-rengan/pkg/tracing"
	mockpublisher "go-rengan/todo/mocks/publisher"
	mockrepository "go-rengan/todo/mocks/repository"
	"go-rengan/todo/models"
	"go-rengan/todo/service"
	errorsutil "go-rengan/utils/errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var DefaultID string = "1"

func TestGetAll(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	t.Run("success when find all", func(t *testing.T) {
		mockList := make([]*models.Todo, 0)
		mockList = append(mockList, &models.Todo{})

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("FindAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(mockList, nil)
		mockRepository.On("CountFindAll", mock.Anything, mock.AnythingOfType("string")).Return(10, nil)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		results, count, err := service.GetAll(context.Background(), "keyword", 10, 0)

		assert.NoError(t, err)
		assert.Equal(t, count, 10)
		assert.Equal(t, mockList, results)
	})

	t.Run("error when find all", func(t *testing.T) {
		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("FindAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, errorsutil.ErrDefault)
		mockRepository.On("CountFindAll", mock.Anything, mock.AnythingOfType("string")).Return(10, nil)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		results, count, err := service.GetAll(context.Background(), "keyword", 10, 0)

		assert.Nil(t, results)
		assert.Equal(t, 0, count)
		assert.Error(t, err)
	})

	t.Run("error when count find all", func(t *testing.T) {
		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("FindAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, nil)
		mockRepository.On("CountFindAll", mock.Anything, mock.AnythingOfType("string")).Return(10, errorsutil.ErrDefault)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		results, count, err := service.GetAll(context.Background(), "keyword", 10, 0)

		assert.Nil(t, results)
		assert.Equal(t, 0, count)
		assert.Error(t, err)
	})
}

func TestGetByID(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	t.Run("success when find by id", func(t *testing.T) {
		var mockTodo = &models.Todo{}

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("FindById", mock.Anything, mock.AnythingOfType("string")).Return(mockTodo, nil)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		result, err := service.GetByID(context.Background(), DefaultID)

		assert.NoError(t, err)
		assert.Equal(t, mockTodo, result)
	})

	t.Run("error when find by id", func(t *testing.T) {
		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("FindById", mock.Anything, mock.AnythingOfType("string")).Return(nil, errorsutil.ErrDefault)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		result, err := service.GetByID(context.Background(), DefaultID)

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestCreate(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	t.Run("success when create", func(t *testing.T) {
		var mockTodo = &models.Todo{}

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("Store", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(mockTodo, nil)

		mockPublisher := new(mockpublisher.AMQPPublisher)
		mockPublisher.On("Create", mock.AnythingOfType("string"))

		service := service.New(tracing, mockRepository, mockPublisher)

		result, err := service.Create(context.Background(), &models.Todo{})

		assert.NoError(t, err)
		assert.Equal(t, mockTodo, result)
	})

	t.Run("error when create", func(t *testing.T) {
		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("Store", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(nil, errorsutil.ErrDefault)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		result, err := service.Create(context.Background(), &models.Todo{})

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestUpdate(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	t.Run("success when update", func(t *testing.T) {
		var mockTodo = &models.Todo{}

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("CountFindByID", mock.Anything, mock.AnythingOfType("string")).Return(10, nil)
		mockRepository.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(mockTodo, nil)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		result, err := service.Update(context.Background(), DefaultID, &models.Todo{})

		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("error when count find by id", func(t *testing.T) {
		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("CountFindByID", mock.Anything, mock.AnythingOfType("string")).Return(0, errorsutil.ErrDefault)
		mockRepository.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(nil, nil)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		result, err := service.Update(context.Background(), DefaultID, &models.Todo{})

		assert.Nil(t, result)
		assert.Error(t, err)
	})

	t.Run("error when update", func(t *testing.T) {
		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("CountFindByID", mock.Anything, mock.AnythingOfType("string")).Return(10, nil)
		mockRepository.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(nil, errorsutil.ErrDefault)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		result, err := service.Update(context.Background(), DefaultID, &models.Todo{})

		assert.Nil(t, result)
		assert.Error(t, err)
	})
}

func TestDelete(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	t.Run("success when delete", func(t *testing.T) {
		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		err = service.Delete(context.Background(), DefaultID)

		assert.NoError(t, err)
	})

	t.Run("error when delete", func(t *testing.T) {
		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockRepository := new(mockrepository.Repository)
		mockRepository.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(errorsutil.ErrDefault)

		mockPublisher := new(mockpublisher.AMQPPublisher)

		service := service.New(tracing, mockRepository, mockPublisher)

		err = service.Delete(context.Background(), DefaultID)

		assert.Error(t, err)
	})
}
