package service

import (
	"context"
	tracing "go-rengan/pkg/tracing"
	"go-rengan/todo/models"
	amqpservice "go-rengan/todo/publisher"
	"go-rengan/todo/repository"
)

// Service represent the todo service
type Service interface {
	GetAll(ctx context.Context, keyword string, limit int, offset int) ([]*models.Todo, int, error)
	GetByID(ctx context.Context, id string) (*models.Todo, error)
	Create(ctx context.Context, value *models.Todo) (*models.Todo, error)
	Update(ctx context.Context, id string, value *models.Todo) (*models.Todo, error)
	Delete(ctx context.Context, id string) error
}

type ServiceImpl struct {
	tracing           tracing.Tracing
	todoRepo          repository.Repository
	todoAMQPPublisher amqpservice.AMQPPublisher
}

// New will create new an ServiceImpl object representation of Service interface
func New(
	tracing tracing.Tracing,
	todoRepo repository.Repository,
	todoAMQPPublisher amqpservice.AMQPPublisher,
) Service {
	return &ServiceImpl{
		tracing:           tracing,
		todoRepo:          todoRepo,
		todoAMQPPublisher: todoAMQPPublisher,
	}
}

// GetAll - get all todo service
func (s *ServiceImpl) GetAll(ctx context.Context, keyword string, limit int, offset int) ([]*models.Todo, int, error) {
	ctx, span := s.tracing.Tracer("TodoService").Start(ctx, "TodoService.GetAll")
	defer span.End()

	res, err := s.todoRepo.FindAll(ctx, keyword, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// Count total
	total, err := s.todoRepo.CountFindAll(ctx, keyword)
	if err != nil {
		return nil, 0, err
	}

	return res, total, nil
}

// GetByID - get todo by id service
func (s *ServiceImpl) GetByID(ctx context.Context, id string) (*models.Todo, error) {
	ctx, span := s.tracing.Tracer("TodoService").Start(ctx, "TodoService.GetByID")
	defer span.End()

	res, err := s.todoRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// Create - creating todo service
func (s *ServiceImpl) Create(ctx context.Context, value *models.Todo) (*models.Todo, error) {
	ctx, span := s.tracing.Tracer("TodoService").Start(ctx, "TodoService.Create")
	defer span.End()

	res, err := s.todoRepo.Store(ctx, &models.Todo{
		Title:       value.Title,
		Description: value.Description,
	})
	if err != nil {
		return nil, err
	}

	// Send Email Queue
	s.todoAMQPPublisher.Create("example.com")

	return res, nil
}

// Update - update todo service
func (s *ServiceImpl) Update(ctx context.Context, id string, value *models.Todo) (*models.Todo, error) {
	ctx, span := s.tracing.Tracer("TodoService").Start(ctx, "TodoService.Update")
	defer span.End()

	_, err := s.todoRepo.CountFindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_, err = s.todoRepo.Update(ctx, id, &models.Todo{
		Title:       value.Title,
		Description: value.Description,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Delete - delete todo service
func (s *ServiceImpl) Delete(ctx context.Context, id string) error {
	ctx, span := s.tracing.Tracer("TodoService").Start(ctx, "TodoService.Delete")
	defer span.End()

	err := s.todoRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
