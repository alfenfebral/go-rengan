package service

import (
	"context"
	pkg_tracing "go-rengan/pkg/tracing"
	"go-rengan/todo/models"
	todo_amqp_service "go-rengan/todo/publisher"
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
	tracing           pkg_tracing.Tracing
	todoRepo          repository.MongoTodoRepository
	todoAMQPPublisher todo_amqp_service.TodoAMQPPublisher
}

// NewTodoService will create new an TodoService object representation of TodoService interface
func NewTodoService(
	tracing pkg_tracing.Tracing,
	todoRepo repository.MongoTodoRepository,
	todoAMQPPublisher todo_amqp_service.TodoAMQPPublisher,
) TodoService {
	return &TodoServiceImpl{
		tracing:           tracing,
		todoRepo:          todoRepo,
		todoAMQPPublisher: todoAMQPPublisher,
	}
}

// GetAll - get all todo service
func (serviceImpl *TodoServiceImpl) GetAll(ctx context.Context, keyword string, limit int, offset int) ([]*models.Todo, int, error) {
	ctx, span := serviceImpl.tracing.Tracer("TodoService").Start(ctx, "TodoService.GetAll")
	defer span.End()

	res, err := serviceImpl.todoRepo.FindAll(ctx, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Count total
	total, err := serviceImpl.todoRepo.CountFindAll(ctx, keyword)
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

// GetByID - get todo by id service
func (serviceImpl *TodoServiceImpl) GetByID(ctx context.Context, id string) (*models.Todo, error) {
	ctx, span := serviceImpl.tracing.Tracer("TodoService").Start(ctx, "TodoService.GetByID")
	defer span.End()

	res, err := serviceImpl.todoRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create - creating todo service
func (serviceImpl *TodoServiceImpl) Create(ctx context.Context, value *models.Todo) (*models.Todo, error) {
	ctx, span := serviceImpl.tracing.Tracer("TodoService").Start(ctx, "TodoService.Create")
	defer span.End()

	res, err := serviceImpl.todoRepo.Store(ctx, &models.Todo{
		Title:       value.Title,
		Description: value.Description,
	})
	if err != nil {
		return nil, err
	}

	// Send Email Queue
	serviceImpl.todoAMQPPublisher.Create("example.com")

	return res, nil
}

// Update - update todo service
func (serviceImpl *TodoServiceImpl) Update(ctx context.Context, id string, value *models.Todo) (*models.Todo, error) {
	ctx, span := serviceImpl.tracing.Tracer("TodoService").Start(ctx, "TodoService.Update")
	defer span.End()

	_, err := serviceImpl.todoRepo.CountFindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = serviceImpl.todoRepo.Update(ctx, id, &models.Todo{
		Title:       value.Title,
		Description: value.Description,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete - delete todo service
func (serviceImpl *TodoServiceImpl) Delete(ctx context.Context, id string) error {
	ctx, span := serviceImpl.tracing.Tracer("TodoService").Start(ctx, "TodoService.Delete")
	defer span.End()

	err := serviceImpl.todoRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
