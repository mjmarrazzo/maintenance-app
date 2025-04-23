package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjmarrazzo/maintenance-app/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/hashing"
	"github.com/mjmarrazzo/maintenance-app/repository"
)

type UserService interface {
	Create(ctx context.Context, user *domain.UserRequest) error
	Authenticate(ctx context.Context, email, password string) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(pool *pgxpool.Pool) UserService {
	return &userService{repo: repository.NewUserRepository(pool)}
}

func (s *userService) Create(ctx context.Context, userRequest *domain.UserRequest) error {
	user := userRequest.ToDomain()

	passwordHash, err := hashing.HashPassword(userRequest.Password)
	if err != nil {
		return err
	}

	user.PasswordHash = passwordHash

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (s *userService) Authenticate(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if ok := hashing.VerifyPassword(password, user.PasswordHash); !ok {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}
