package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mjmarrazzo/maintenance-app/internal/database"
	"github.com/mjmarrazzo/maintenance-app/internal/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/responses"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id int64) (*domain.Task, error)
	GetAll(ctx context.Context, filters TaskFilters) ([]*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id int64) error
	UpdateStatus(ctx context.Context, id int64, status domain.Status) error
	AssignTask(ctx context.Context, taskID int64, userID int64) error
	CompleteTask(ctx context.Context, id int64) error
	CountByStatus(ctx context.Context) (map[domain.Status]int, error)
	CountByPriority(ctx context.Context) (map[domain.Priority]int, error)
}

type TaskFilters struct {
	Status      *domain.Status
	Priority    *domain.Priority
	CategoryID  *int64
	LocationID  *int64
	AssignedTo  *int64
	CreatedBy   *int64
	IsCompleted *bool
	IsRecurring *bool
	SearchQuery string
	DateFrom    *time.Time
	DateTo      *time.Time
	Limit       int
	Offset      int
	SortField   string
	SortOrder   string
}

type taskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	query := `
		INSERT INTO tasks (
			title,
			description,
			category_id,
			location_id,
			priority,
			status,
			created_by,
			assigned_to,
			estimated_completion_date,
			cost,
			is_recurring,
			recurrence_type,
			recurrence_interval,
			recurrence_unit,
			parent_task_id
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14, $15
		) RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.CategoryID,
		task.LocationID,
		task.Priority,
		task.Status,
		task.CreatedBy,
		task.AssignedTo,
		task.EstimatedCompletionDate,
		task.Cost,
		task.IsRecurring,
		task.RecurrenceType,
		task.RecurrenceInterval,
		task.RecurrenceUnit,
		task.ParentTaskID,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		return database.HandleError(err, "task", task.ID)
	}

	return nil
}

func scanRowToTask(row pgx.Row, task *domain.Task) error {
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.CategoryID,
		&task.CategoryName,
		&task.LocationID,
		&task.LocationName,
		&task.Priority,
		&task.Status,
		&task.CreatedBy,
		&task.CreatedByName,
		&task.AssignedTo,
		&task.AssignedToName,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.EstimatedCompletionDate,
		&task.Cost,
		&task.IsRecurring,
		&task.RecurrenceType,
		&task.RecurrenceInterval,
		&task.RecurrenceUnit,
		&task.ParentTaskID,
		&task.NextOccurrence,
	)
	if err != nil {
		return fmt.Errorf("error scanning task: %w", err)
	}
	return nil
}
func (r *taskRepository) GetByID(ctx context.Context, id int64) (*domain.Task, error) {
	query := `
		SELECT
			t.id,
			t.title,
			t.description,
			t.category_id,
			c.name AS category_name,
			t.location_id,
			l.name AS location_name,
			t.priority,
			t.status,
			t.created_by,
			creator.name AS creator_name,
			t.assigned_to,
			assignee.name AS assigned_to_name,
			t.created_at,
			t.updated_at,
			t.estimated_completion_date,
			t.cost,
			t.is_recurring,
			t.recurrence_type,
			t.recurrence_interval,
			t.recurrence_unit,
			t.parent_task_id,
			t.next_occurrence
		FROM tasks t
		LEFT JOIN categories c ON t.category_id = c.id
		LEFT JOIN locations l ON t.location_id = l.id
		LEFT JOIN users creator ON t.created_by = creator.id
		LEFT JOIN users assignee ON t.assigned_to = assignee.id
		WHERE t.id = $1;
	`

	row := r.db.QueryRow(ctx, query, id)

	task := &domain.Task{}
	if err := scanRowToTask(row, task); err != nil {
		return nil, database.HandleError(err, "task", id)
	}

	return task, nil
}

func (r *taskRepository) GetAll(ctx context.Context, filters TaskFilters) ([]*domain.Task, error) {
	query := `
		SELECT
			t.id,
			t.title,
			t.description,
			t.category_id,
			c.name AS category_name,
			t.location_id,
			l.name AS location_name,
			t.priority,
			t.status,
			t.created_by,
			creator.name AS creator_name,
			t.assigned_to,
			assignee.name AS assigned_to_name,
			t.created_at,
			t.updated_at,
			t.estimated_completion_date,
			t.cost,
			t.is_recurring,
			t.recurrence_type,
			t.recurrence_interval,
			t.recurrence_unit,
			t.parent_task_id,
			t.next_occurrence
		FROM tasks t
		LEFT JOIN categories c ON t.category_id = c.id
		LEFT JOIN locations l ON t.location_id = l.id
		LEFT JOIN users creator ON t.created_by = creator.id
		LEFT JOIN users assignee ON t.assigned_to = assignee.id
		WHERE 1=1
	`
	var args []interface{}
	argIndex := 1

	if filters.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filters.Status)
		argIndex++
	}

	if filters.Priority != nil {
		query += fmt.Sprintf(" AND priority = $%d", argIndex)
		args = append(args, *filters.Priority)
		argIndex++
	}

	if filters.CategoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, *filters.CategoryID)
		argIndex++
	}

	if filters.LocationID != nil {
		query += fmt.Sprintf(" AND location_id = $%d", argIndex)
		args = append(args, *filters.LocationID)
		argIndex++
	}

	if filters.AssignedTo != nil {
		query += fmt.Sprintf(" AND assigned_to = $%d", argIndex)
		args = append(args, *filters.AssignedTo)
		argIndex++
	}

	if filters.CreatedBy != nil {
		query += fmt.Sprintf(" AND created_by = $%d", argIndex)
		args = append(args, *filters.CreatedBy)
		argIndex++
	}

	if filters.IsCompleted != nil {
		if *filters.IsCompleted {
			query += " AND status = 'Completed'"
		} else {
			query += " AND status != 'Completed'"
		}
	}

	if filters.IsRecurring != nil {
		query += fmt.Sprintf(" AND is_recurring = $%d", argIndex)
		args = append(args, *filters.IsRecurring)
		argIndex++
	}

	if filters.SearchQuery != "" {
		query += fmt.Sprintf(" AND (title ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex)
		searchPattern := "%" + filters.SearchQuery + "%"
		args = append(args, searchPattern)
		argIndex++
	}

	if filters.DateFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argIndex)
		args = append(args, *filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argIndex)
		args = append(args, *filters.DateTo)
		argIndex++
	}

	if filters.SortField != "" {
		sortOrder := "ASC"
		if strings.ToUpper(filters.SortOrder) == "DESC" {
			sortOrder = "DESC"
		}

		allowedFields := map[string]bool{
			"id": true, "title": true, "priority": true, "status": true,
			"created_at": true, "updated_at": true, "estimated_completion_date": true,
		}

		if allowedFields[filters.SortField] {
			query += fmt.Sprintf(" ORDER BY %s %s", filters.SortField, sortOrder)
		} else {
			query += " ORDER BY created_at DESC"
		}
	} else {
		query += " ORDER BY created_at DESC"
	}

	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error listing tasks: %w", err)
	}
	defer rows.Close()

	tasks := []*domain.Task{}
	for rows.Next() {
		task := &domain.Task{}
		if err := scanRowToTask(rows, task); err != nil {
			return nil, database.HandleError(err, "task", filters)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `
		UPDATE tasks SET
			title = $1,
			description = $2,
			category_id = $3,
			location_id = $4,
			priority = $5,
			status = $6,
			assigned_to = $7,
			estimated_completion_date = $8,
			cost = $9,
			is_recurring = $10,
			recurrence_type = $11,
			recurrence_interval = $12,
			recurrence_unit = $13,
			parent_task_id = $14,
			updated_at = NOW()
		WHERE id = $15
		RETURNING updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.CategoryID,
		task.LocationID,
		task.Priority,
		task.Status,
		task.AssignedTo,
		task.EstimatedCompletionDate,
		task.Cost,
		task.IsRecurring,
		task.RecurrenceType,
		task.RecurrenceInterval,
		task.RecurrenceUnit,
		task.ParentTaskID,
		task.ID,
	).Scan(&task.UpdatedAt)

	if err != nil {
		return database.HandleError(err, "task", task.ID)
	}

	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM tasks WHERE id = $1`
	ct, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return database.HandleError(err, "task", id)
	}

	if ct.RowsAffected() == 0 {
		return responses.NewNotFoundError(fmt.Sprintf("task with ID %d not found", id))
	}

	return nil
}

func (r *taskRepository) UpdateStatus(ctx context.Context, id int64, status domain.Status) error {
	query := `
		UPDATE tasks SET
			status = $1,
			completed_at = CASE WHEN $1 = 'Completed' THEN NOW() ELSE NULL END
		WHERE id = $2
		RETURNING updated_at`

	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, status, id).Scan(&updatedAt)
	if err != nil {
		return database.HandleError(err, "task", id)
	}

	return nil
}

func (r *taskRepository) AssignTask(ctx context.Context, taskID int64, userID int64) error {
	query := `UPDATE tasks SET assigned_to = $1 WHERE id = $2 RETURNING updated_at`

	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, userID, taskID).Scan(&updatedAt)
	if err != nil {
		return database.HandleError(err, "task", taskID)

	}

	return nil
}

func (r *taskRepository) CompleteTask(ctx context.Context, id int64) error {
	query := `
		UPDATE tasks SET
			status = 'Completed',
			completed_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&updatedAt)
	if err != nil {
		return database.HandleError(err, "task", id)

	}

	return nil
}

func (r *taskRepository) CountByStatus(ctx context.Context) (map[domain.Status]int, error) {
	query := `
		SELECT status, COUNT(*) as count
		FROM tasks
		GROUP BY status`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error counting tasks by status: %w", err)
	}
	defer rows.Close()

	result := make(map[domain.Status]int)
	for rows.Next() {
		var status domain.Status
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("error scanning status count: %w", err)
		}
		result[status] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating status counts: %w", err)
	}

	return result, nil
}

func (r *taskRepository) CountByPriority(ctx context.Context) (map[domain.Priority]int, error) {
	query := `
		SELECT priority, COUNT(*) as count
		FROM tasks
		GROUP BY priority`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error counting tasks by priority: %w", err)
	}
	defer rows.Close()

	result := make(map[domain.Priority]int)
	for rows.Next() {
		var priority domain.Priority
		var count int
		if err := rows.Scan(&priority, &count); err != nil {
			return nil, fmt.Errorf("error scanning priority count: %w", err)
		}
		result[priority] = count
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating priority counts: %w", err)
	}

	return result, nil
}
