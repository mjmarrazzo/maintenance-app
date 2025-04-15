package database

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/mjmarrazzo/maintenance-app/internal/responses"
)

func IsPgError(err error, code string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == code
}

func IsUniqueViolation(err error) bool {
	return IsPgError(err, UniqueViolation)
}

func IsForeignKeyViolation(err error) bool {
	return IsPgError(err, ForeignKeyViolation)
}

func IsNotNullViolation(err error) bool {
	return IsPgError(err, NotNullViolation)
}

func IsNotFoundViolation(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}

func HandleError(err error, entityType string, id interface{}) error {
	if IsUniqueViolation(err) {
		return responses.NewConflictError(
			fmt.Sprintf("%s with ID %v already exists", entityType, id),
		)
	} else if IsForeignKeyViolation(err) {
		fmt.Printf("Foreign key violation: %v\n", err)
		return &Error{
			Err:        ErrNotFound,
			EntityType: entityType,
			EntityID:   id,
			Op:         "delete",
		}
	} else if IsNotNullViolation(err) {
		fmt.Printf("Not null violation: %v\n", err)
		return &Error{
			Err:        ErrNotFound,
			EntityType: entityType,
			EntityID:   id,
			Op:         "update",
		}
	} else if IsNotFoundViolation(err) {
		fmt.Printf("Not found violation: %v\n", err)
		return responses.NewNotFoundError(
			fmt.Sprintf("%s with ID %v not found", entityType, id),
		)
	}

	return err
}
