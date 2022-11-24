package httpdelivery_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	tracing "go-rengan/pkg/tracing"
	validator "go-rengan/pkg/validator"
	errorsutil "go-rengan/utils/errors"

	httpdelivery "go-rengan/todo/delivery/http"
	mockservice "go-rengan/todo/mocks/service"

	"go-rengan/todo/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var WhenError400EOF string = "when return 400 bad request (error EOF)"
var WhenError500Service string = "when return 500 internal error (error service)"
var WhenError500Query string = "when return 500 internal error (error query)"
var WhenError400Validation string = "when return 400 bad request (error validation)"
var WhenError404NotFound string = "when return 404 not found (resouce not found)"
var WhenSuccess201Created string = "when return 201 created"
var WhenSuccess200OK string = "when return 200 ok"

func TestNew(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	validator.New()

	tracing, err := tracing.New()
	assert.NoError(t, err)

	mockservice := new(mockservice.Service)
	mockservice.On("Create", mock.AnythingOfType("*models.Todo")).Return(&models.Todo{}, nil)

	handler := httpdelivery.New(tracing, mockservice)
	router := chi.NewMux()
	handler.RegisterRoutes(router)
}

// TestGetAll - testing GetAll [200]
func TestGetAll(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	t.Run(WhenError400Validation, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?page=-1&per_page=-1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetAll)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})

	t.Run(WhenError500Service, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?page=1&per_page=10", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(nil, 1, errorsutil.ErrDefault)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetAll)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})

	t.Run(WhenSuccess200OK, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?page=1&per_page=10", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		mockListTodo := make([]*models.Todo, 0)
		mockListTodo = append(mockListTodo, &models.Todo{})

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("GetAll", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("int"), mock.AnythingOfType("int")).Return(mockListTodo, 1, nil)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetAll)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
}

// TestCreate - testing create [201]
func TestCreate(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	t.Run(WhenError400EOF, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader([]byte("")))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("when return 400 bad request (error validation) ", func(t *testing.T) {
		validator.New()

		mockPostBody := map[string]interface{}{
			"title":       "",
			"description": "",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("when error 500 internal error (error service)", func(t *testing.T) {
		validator.New()

		mockPostBody := map[string]interface{}{
			"title":       "lorem ipsum",
			"description": "desc",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("Create", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(nil, errorsutil.ErrDefault)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})

	t.Run("when return 201 created", func(t *testing.T) {
		validator.New()

		mockPostBody := map[string]interface{}{
			"title":       "lorem ipsum",
			"description": "desc",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/todo", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("Create", mock.Anything, mock.AnythingOfType("*models.Todo")).Return(&models.Todo{}, nil)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Create)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusCreated, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
}

// TestGetByID - testing GetByID [200]
func TestGetByID(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	t.Run(WhenError404NotFound, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errorsutil.ErrNotFound)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetByID)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusNotFound, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(nil, errorsutil.ErrDefault)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetByID)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodGet, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("GetByID", mock.Anything, mock.AnythingOfType("string")).Return(&models.Todo{}, nil)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.GetByID)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
}

// TestUpdate - testing update [200]
func TestUpdate(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	t.Run(WhenError400EOF, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodPut, "/api/v1/product?id=1", bytes.NewReader([]byte("")))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(&models.Todo{}, nil)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})
	t.Run(WhenError400Validation, func(t *testing.T) {
		validator.New()

		mockPostBody := map[string]interface{}{
			"title":       "",
			"description": "",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusBadRequest, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
	t.Run(WhenError404NotFound, func(t *testing.T) {
		validator.New()

		mockPostBody := map[string]interface{}{
			"title":       "a",
			"description": "a",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(nil, errorsutil.ErrNotFound)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusNotFound, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		validator.New()

		mockPostBody := map[string]interface{}{
			"title":       "a",
			"description": "a",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(nil, errorsutil.ErrDefault)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		validator.New()

		mockPostBody := map[string]interface{}{
			"title":       "a",
			"description": "a",
		}
		body, _ := json.Marshal(mockPostBody)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/todo?id=1", bytes.NewReader(body))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("Update", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("*models.Todo")).Return(&models.Todo{}, nil)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Update)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
}

// TestDeleteSuccess - testing delete [200]
func TestDelete(t *testing.T) {
	os.Setenv("APP_ID", "1")
	os.Setenv("APP_NAME", "go-rengan")
	os.Setenv("TRACER_PROVIDER_URL", "http://project2_secret_token@localhost:14317/2")

	t.Run(WhenError404NotFound, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/product?id=", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(errorsutil.ErrNotFound)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Delete)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})
	t.Run(WhenError500Service, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(errorsutil.ErrDefault)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Delete)
		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
	t.Run(WhenSuccess200OK, func(t *testing.T) {
		validator.New()

		req, err := http.NewRequest(http.MethodDelete, "/api/v1/todo?id=1", nil)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/json")

		tracing, err := tracing.New()
		assert.NoError(t, err)

		mockservice := new(mockservice.Service)
		mockservice.On("Delete", mock.Anything, mock.AnythingOfType("string")).Return(nil)

		todoHandler := httpdelivery.New(tracing, mockservice)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(todoHandler.Delete)

		handler.ServeHTTP(rr, req)

		// Check the status code is what expected
		assert.Equal(t, http.StatusOK, rr.Code)

		// Check if the mock called
		mockservice.AssertExpectations(t)
	})
}
