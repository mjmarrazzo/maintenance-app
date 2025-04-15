package domain

import (
	"database/sql"
	"strconv"
	"time"
)

type Priority string

const (
	PriorityLow    Priority = "Low"
	PriorityMedium Priority = "Medium"
	PriorityHigh   Priority = "High"
	PriorityUrgent Priority = "Urgent"
)

type Status string

const (
	StatusNew        Status = "New"
	StatusInProgress Status = "In Progress"
	StatusCompleted  Status = "Completed"
	StatusOnHold     Status = "On Hold"
)

type RecurrenceType string

const (
	RecurrenceTypeDaily   RecurrenceType = "Daily"
	RecurrenceTypeWeekly  RecurrenceType = "Weekly"
	RecurrenceTypeMonthly RecurrenceType = "Monthly"
	RecurrenceTypeYearly  RecurrenceType = "Yearly"
	RecurrentTypeCustom   RecurrenceType = "Custom"
)

type RecurrenceUnit string

const (
	RecurrenceUnitDay   RecurrenceUnit = "Day"
	RecurrenceUnitWeek  RecurrenceUnit = "Week"
	RecurrenceUnitMonth RecurrenceUnit = "Month"
	RecurrenceUnitYear  RecurrenceUnit = "Year"
)

type Task struct {
	ID                      int64         `db:"id"`
	Title                   string        `db:"title"`
	Description             string        `db:"description"`
	CategoryID              sql.NullInt64 `db:"category_id"`
	CategoryName            sql.NullString
	LocationID              sql.NullInt64 `db:"location_id"`
	LocationName            sql.NullString
	Priority                sql.NullString `db:"priority"`
	Status                  sql.NullString `db:"status"`
	CreatedBy               int64          `db:"created_by"`
	CreatedByName           sql.NullString
	AssignedTo              sql.NullInt64 `db:"assigned_to"`
	AssignedToName          sql.NullString
	CreatedAt               time.Time       `db:"created_at"`
	UpdatedAt               time.Time       `db:"updated_at"`
	EstimatedCompletionDate sql.NullTime    `db:"estimated_completion_date"`
	Cost                    sql.NullFloat64 `db:"cost"`
	IsRecurring             bool            `db:"is_recurring"`
	RecurrenceType          sql.NullString  `db:"recurrence_type"`
	RecurrenceInterval      int             `db:"recurrence_interval"`
	RecurrenceUnit          sql.NullString  `db:"recurrence_unit"`
	ParentTaskID            sql.NullInt64   `db:"parent_task_id"`
	NextOccurrence          sql.NullTime    `db:"next_occurrence"`
	CompletedAt             sql.NullTime    `db:"completed_at"`
}

type TaskRequest struct {
	Title                   string `json:"title" form:"title"`
	Description             string `json:"description" form:"description"`
	CategoryID              string `json:"category_id" form:"category_id"`
	LocationID              string `json:"location_id" form:"location_id"`
	TaskPriority            string `json:"priority" form:"priority"`
	TaskStatus              string `json:"status" form:"status"`
	AssignedTo              string `json:"assigned_to" form:"assigned_to"`
	EstimatedCompletionDate string `json:"estimated_completion_date" form:"estimated_completion_date"`
	Cost                    string `json:"cost" form:"cost"`
	IsRecurring             bool   `json:"is_recurring" form:"is_recurring"`
	RecurrenceType          string `json:"recurrence_type" form:"recurrence_type"`
	RecurrenceInterval      int    `json:"recurrence_interval" form:"recurrence_interval"`
	RecurrenceUnit          string `json:"recurrence_unit" form:"recurrence_unit"`
	ParentTaskID            string `json:"parent_task_id" form:"parent_task_id"`
}

func (tr *TaskRequest) ToDomain() *Task {
	var categoryID, locationID, assignedTo, parentTaskID sql.NullInt64
	if tr.CategoryID != "" {
		num, err := strconv.ParseInt(tr.CategoryID, 10, 64)
		if err == nil {
			categoryID.Int64 = num
			categoryID.Valid = true
		}
	}

	if tr.LocationID != "" {
		num, err := strconv.ParseInt(tr.LocationID, 10, 64)
		if err == nil {
			locationID.Int64 = num
			locationID.Valid = true
		}
	}

	if tr.AssignedTo != "" {
		num, err := strconv.ParseInt(tr.AssignedTo, 10, 64)
		if err == nil {
			assignedTo.Int64 = num
			assignedTo.Valid = true
		}
	}

	if tr.ParentTaskID != "" {
		num, err := strconv.ParseInt(tr.ParentTaskID, 10, 64)
		if err == nil {
			parentTaskID.Int64 = num
			parentTaskID.Valid = true
		}
	}

	var cost sql.NullFloat64
	if tr.Cost != "" {
		num, err := strconv.ParseFloat(tr.Cost, 64)
		if err == nil {
			cost.Float64 = num
			cost.Valid = true
		}
	}

	var estimatedCompletionDate sql.NullTime
	if tr.EstimatedCompletionDate != "" {
		t, err := time.Parse(time.RFC3339, tr.EstimatedCompletionDate)
		if err == nil {
			estimatedCompletionDate.Time = t
			estimatedCompletionDate.Valid = true
		}
	}

	var recurrenceType sql.NullString
	if tr.RecurrenceType != "" {
		recurrenceType.String = string(RecurrenceType(tr.RecurrenceType))
		recurrenceType.Valid = true
	}

	var recurrenceUnit sql.NullString
	if tr.RecurrenceUnit != "" {
		recurrenceUnit.String = string(RecurrenceUnit(tr.RecurrenceUnit))
		recurrenceUnit.Valid = true
	}

	var priority sql.NullString
	if tr.TaskPriority != "" {
		priority.String = string(Priority(tr.TaskPriority))
		priority.Valid = true
	}

	var status sql.NullString
	if tr.TaskStatus != "" {
		status.String = string(Status(tr.TaskStatus))
		status.Valid = true
	}

	return &Task{
		Title:                   tr.Title,
		Description:             tr.Description,
		CategoryID:              categoryID,
		LocationID:              locationID,
		Priority:                priority,
		Status:                  status,
		CreatedBy:               1, //TODO: actually fill this in
		AssignedTo:              assignedTo,
		EstimatedCompletionDate: estimatedCompletionDate,
		Cost:                    cost,
		IsRecurring:             tr.IsRecurring,
		RecurrenceType:          recurrenceType,
		RecurrenceInterval:      tr.RecurrenceInterval,
		RecurrenceUnit:          recurrenceUnit,
		ParentTaskID:            parentTaskID,
	}
}
