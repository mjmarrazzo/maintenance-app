package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjmarrazzo/maintenance-app/domain"
	"github.com/mjmarrazzo/maintenance-app/repository"
)

type CategoryService interface {
	Create(ctx context.Context, category *domain.CategoryRequest) (*domain.Category, error)
	GetAll(ctx context.Context) ([]*domain.Category, error)
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	Update(ctx context.Context, id int64, category *domain.CategoryRequest) (*domain.Category, error)
	Delete(ctx context.Context, id int64) error
}

type categoryService struct {
	repository repository.CategoryRepository
}

func NewCategoryService(pool *pgxpool.Pool) CategoryService {
	return &categoryService{repository: repository.NewCategoryRepository(pool)}
}

func (s *categoryService) Create(ctx context.Context, category *domain.CategoryRequest) (*domain.Category, error) {
	categoryDomain := category.ToDomain()
	if err := s.repository.Create(ctx, categoryDomain); err != nil {
		return nil, err
	}
	return categoryDomain, nil
}

func (s *categoryService) GetAll(ctx context.Context) ([]*domain.Category, error) {
	return s.repository.GetAll(ctx)
}

func (s *categoryService) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *categoryService) Update(ctx context.Context, id int64, category *domain.CategoryRequest) (*domain.Category, error) {
	categoryDomain := category.ToDomain()
	categoryDomain.ID = id
	if err := s.repository.Update(ctx, categoryDomain); err != nil {
		return nil, err
	}
	return categoryDomain, nil
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
