package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjmarrazzo/maintenance-app/internal/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/repository"
)

type LocationService interface {
	Create(ctx context.Context, location *domain.LocationRequest) (*domain.Location, error)
	GetAll(ctx context.Context) ([]*domain.Location, error)
	GetByID(ctx context.Context, id int64) (*domain.Location, error)
	Update(ctx context.Context, id int64, location *domain.LocationRequest) (*domain.Location, error)
	Delete(ctx context.Context, id int64) error
}

type locationService struct {
	repository repository.LocationRepository
}

func NewLocationService(pool *pgxpool.Pool) LocationService {
	return &locationService{repository: repository.NewLocationRepository(pool)}
}

func (s *locationService) Create(ctx context.Context, location *domain.LocationRequest) (*domain.Location, error) {
	locationDomain := location.ToDomain()
	if err := s.repository.Create(ctx, locationDomain); err != nil {
		return nil, err
	}
	return locationDomain, nil
}

func (s *locationService) GetAll(ctx context.Context) ([]*domain.Location, error) {
	return s.repository.GetAll(ctx)
}

func (s *locationService) GetByID(ctx context.Context, id int64) (*domain.Location, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *locationService) Update(ctx context.Context, id int64, location *domain.LocationRequest) (*domain.Location, error) {
	locationDomain := location.ToDomain()
	locationDomain.ID = id
	if err := s.repository.Update(ctx, locationDomain); err != nil {
		return nil, err
	}
	return locationDomain, nil
}

func (s *locationService) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
