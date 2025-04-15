package database

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type Error struct {
	Err        error
	EntityType string
	EntityID   interface{}
	Op         string
}

func (e *Error) Error() string {
	if e.EntityID != nil {
		return fmt.Sprintf("%s: %s %v with ID %v", e.Op, e.Err.Error(), e.EntityType, e.EntityID)
	}
	return fmt.Sprintf("%s: %s %v", e.Op, e.Err.Error(), e.EntityType)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NotFound(entityType string, id interface{}) error {
	return &Error{
		Err:        ErrNotFound,
		EntityType: entityType,
		EntityID:   id,
		Op:         "find",
	}
}

func Conflict(entityType string, id interface{}) error {
	return &Error{
		Err:        ErrConflict,
		EntityType: entityType,
		EntityID:   id,
		Op:         "create",
	}
}
