package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjmarrazzo/maintenance-app/internal/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/repository"
)

type TaskService interface {
	Create(ctx context.Context, userId int64, task *domain.TaskRequest) (*domain.Task, error)
	GetAll(ctx context.Context) ([]*domain.Task, error)
	GetByID(ctx context.Context, id int64) (*domain.Task, error)
	Update(ctx context.Context, id int64, task *domain.TaskRequest) (*domain.Task, error)
	Delete(ctx context.Context, id int64) error
}

type taskService struct {
	repository repository.TaskRepository
}

func NewTaskService(pool *pgxpool.Pool) TaskService {
	return &taskService{repository: repository.NewTaskRepository(pool)}
}

func (s *taskService) Create(ctx context.Context, userId int64, tr *domain.TaskRequest) (*domain.Task, error) {
	task := tr.ToDomain()
	task.CreatedBy = userId

	if err := s.repository.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetAll(ctx context.Context) ([]*domain.Task, error) {
	return s.repository.GetAll(ctx, repository.TaskFilters{})
}

func (s *taskService) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *taskService) Update(ctx context.Context, id int64, tr *domain.TaskRequest) (*domain.Task, error) {
	task := tr.ToDomain()
	task.ID = id

	if err := s.repository.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) Delete(ctx context.Context, id int64) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
