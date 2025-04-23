package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjmarrazzo/maintenance-app/internal/database"
	"github.com/mjmarrazzo/maintenance-app/internal/domain"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	sql := `INSERT INTO users (first_name, last_name, email, password_hash, role) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	row := r.db.QueryRow(ctx, sql,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.Role,
	)

	if err := row.Scan(&user.ID, &user.CreatedAt); err != nil {
		return database.HandleError(err, "user", nil)
	}

	return nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	sql := `SELECT id, first_name, last_name, email, password_hash, role, created_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, sql, email)

	user := &domain.User{}
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt); err != nil {
		return nil, database.HandleError(err, "user", nil)
	}

	return user, nil
}
