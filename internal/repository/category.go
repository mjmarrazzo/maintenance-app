package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjmarrazzo/maintenance-app/internal/database"
	"github.com/mjmarrazzo/maintenance-app/internal/domain"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	GetAll(ctx context.Context) ([]*domain.Category, error)
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int64) error
}

type categoryRepository struct {
	db *pgxpool.Pool
}

func NewCategoryRepository(db *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	sql := `INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id`
	row := r.db.QueryRow(ctx, sql, category.Name, category.Description)

	if err := row.Scan(&category.ID); err != nil {
		return database.HandleError(err, "category", nil)
	}
	return nil
}

func (r *categoryRepository) GetAll(ctx context.Context) ([]*domain.Category, error) {
	sql := `SELECT id, name, description FROM categories`
	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, database.HandleError(err, "category", nil)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		category := &domain.Category{}
		if err := rows.Scan(&category.ID, &category.Name, &category.Description); err != nil {
			return nil, database.HandleError(err, "category", nil)
		}
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, database.HandleError(err, "category", nil)
	}
	return categories, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	sql := `SELECT id, name, description FROM categories WHERE id = $1`
	row := r.db.QueryRow(ctx, sql, id)

	category := &domain.Category{}
	if err := row.Scan(&category.ID, &category.Name, &category.Description); err != nil {
		return nil, database.HandleError(err, "category", id)
	}
	return category, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	sql := `UPDATE categories SET name = $1, description = $2 WHERE id = $3`
	if _, err := r.db.Exec(ctx, sql, category.Name, category.Description, category.ID); err != nil {
		return database.HandleError(err, "category", category.ID)
	}
	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	sql := `DELETE FROM categories WHERE id = $1`
	if _, err := r.db.Exec(ctx, sql, id); err != nil {
		return database.HandleError(err, "category", id)
	}
	return nil
}
