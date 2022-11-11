package service

import (
	"context"
	"go-rengan/todo/models"
	"go-rengan/todo/repository"
)

// TodoService represent the todo service
type TodoService interface {
	GetAll(ctx context.Context, keyword string, limit int, offset int) ([]*models.Todo, int, error)
	GetByID(ctx context.Context, id string) (*models.Todo, error)
	Create(ctx context.Context, value *models.Todo) (*models.Todo, error)
	Update(ctx context.Context, id string, value *models.Todo) (*models.Todo, error)
	Delete(ctx context.Context, id string) error
}

type TodoServiceImpl struct {
	todoRepo repository.MongoTodoRepository
}

// NewTodoService will create new an TodoService object representation of TodoService interface
func NewTodoService(a repository.MongoTodoRepository) TodoService {
	return &TodoServiceImpl{
		todoRepo: a,
	}
}

// GetAll - get all todo service
func (a *TodoServiceImpl) GetAll(ctx context.Context, keyword string, limit int, offset int) ([]*models.Todo, int, error) {
	res, err := a.todoRepo.FindAll(ctx, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Count total
	total, err := a.todoRepo.CountFindAll(ctx, keyword)
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

// GetByID - get todo by id service
func (a *TodoServiceImpl) GetByID(ctx context.Context, id string) (*models.Todo, error) {
	res, err := a.todoRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create - creating todo service
func (a *TodoServiceImpl) Create(ctx context.Context, value *models.Todo) (*models.Todo, error) {
	res, err := a.todoRepo.Store(ctx, &models.Todo{
		Title:       value.Title,
		Description: value.Description,
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Update - update todo service
func (a *TodoServiceImpl) Update(ctx context.Context, id string, value *models.Todo) (*models.Todo, error) {
	_, err := a.todoRepo.CountFindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = a.todoRepo.Update(ctx, id, &models.Todo{
		Title:       value.Title,
		Description: value.Description,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete - delete todo service
func (a *TodoServiceImpl) Delete(ctx context.Context, id string) error {
	err := a.todoRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
