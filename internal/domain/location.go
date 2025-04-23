package domain

import (
	"database/sql"
	"strconv"
)

type Location struct {
	ID                 int64          `db:"id"`
	Name               string         `db:"name"`
	Description        string         `db:"description"`
	ParentLocationId   sql.NullInt64  `db:"parent_location_id"`
	ParentLocationName sql.NullString `db:"parent_location_name"`
}

type LocationRequest struct {
	Name        string `json:"name" form:"name" validate:"required,max=100"`
	Description string `json:"description" form:"description"`
	ParentID    string `json:"parent_location_id" form:"parent_location_id"`
}

func (lr *LocationRequest) ToDomain() *Location {
	var parentID sql.NullInt64
	if lr.ParentID != "" {
		num, err := strconv.ParseInt(lr.ParentID, 10, 64)

		if err != nil {
			parentID.Valid = false
		} else {
			parentID.Int64 = num
			parentID.Valid = true
		}
	}
	return &Location{
		Name:             lr.Name,
		Description:      lr.Description,
		ParentLocationId: parentID,
	}
}
