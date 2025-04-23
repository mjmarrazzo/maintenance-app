package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjmarrazzo/maintenance-app/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/database"
)

type LocationRepository interface {
	Create(ctx context.Context, Location *domain.Location) error
	GetAll(ctx context.Context) ([]*domain.Location, error)
	GetByID(ctx context.Context, id int64) (*domain.Location, error)
	Update(ctx context.Context, Location *domain.Location) error
	Delete(ctx context.Context, id int64) error
}

type locationRepository struct {
	db *pgxpool.Pool
}

func NewLocationRepository(db *pgxpool.Pool) LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(ctx context.Context, location *domain.Location) error {
	sql := `INSERT INTO locations (name, description, parent_location_id) VALUES ($1, $2, $3) RETURNING id`
	row := r.db.QueryRow(ctx, sql, location.Name, location.Description, location.ParentLocationId)

	if err := row.Scan(&location.ID); err != nil {
		return database.HandleError(err, "location", nil)
	}
	return nil
}

func (r *locationRepository) GetAll(ctx context.Context) ([]*domain.Location, error) {
	sql := `
		SELECT l.id, l.name, l.description, l.parent_location_id, p.name AS parent_location_name
		FROM locations l
		LEFT JOIN locations p ON l.parent_location_id = p.id
	`
	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, database.HandleError(err, "location", nil)
	}
	defer rows.Close()

	var locations []*domain.Location
	for rows.Next() {
		location := &domain.Location{}
		if err := rows.Scan(&location.ID, &location.Name, &location.Description, &location.ParentLocationId, &location.ParentLocationName); err != nil {
			return nil, database.HandleError(err, "location", nil)
		}
		locations = append(locations, location)
	}
	if err := rows.Err(); err != nil {
		return nil, database.HandleError(err, "location", nil)
	}
	return locations, nil
}

func (r *locationRepository) GetByID(ctx context.Context, id int64) (*domain.Location, error) {
	sql := `SELECT id, name, description, parent_location_id FROM locations WHERE id = $1`
	row := r.db.QueryRow(ctx, sql, id)

	location := &domain.Location{}
	if err := row.Scan(&location.ID, &location.Name, &location.Description, &location.ParentLocationId); err != nil {
		return nil, database.HandleError(err, "location", nil)
	}
	return location, nil
}

func (r *locationRepository) Update(ctx context.Context, location *domain.Location) error {
	sql := `UPDATE locations SET name = $1, description = $2, parent_location_id = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, sql, location.Name, location.Description, location.ParentLocationId, location.ID)
	if err != nil {
		return database.HandleError(err, "location", location.ID)
	}
	return nil
}

func (r *locationRepository) Delete(ctx context.Context, id int64) error {
	sql := `DELETE FROM locations WHERE id = $1`
	if _, err := r.db.Exec(ctx, sql, id); err != nil {
		return database.HandleError(err, "location", id)
	}
	return nil
}
