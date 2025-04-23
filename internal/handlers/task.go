package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/internal/api"
	"github.com/mjmarrazzo/maintenance-app/internal/auth"
	"github.com/mjmarrazzo/maintenance-app/internal/database"
	"github.com/mjmarrazzo/maintenance-app/internal/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/service"
	"github.com/mjmarrazzo/maintenance-app/internal/views/task_views"
)

type TaskHandler interface {
	api.Handler
	Create(c echo.Context) error
	GetAllTasks(c echo.Context) error
	GetForm(c echo.Context) error
	GetEditForm(c echo.Context) error
	Update(c echo.Context) error
	// Delete(c echo.Context) error
	GetSelect(c echo.Context) error
}

type taskHandler struct {
	service service.TaskService
}

func (c taskHandler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/tasks")
	group.Use(auth.AuthenticatedMiddleware())

	group.POST("", c.Create)
	group.GET("", c.GetAllTasks)
	group.GET("/form", c.GetForm)
	group.GET("/:id/form", c.GetEditForm)
	group.PUT("/:id", c.Update)
	group.DELETE("/:id", c.Delete)
	group.GET("/select", c.GetSelect)
}

func NewTaskHandler(db *database.Client) TaskHandler {
	return &taskHandler{service: service.NewTaskService(db.Pool())}
}

func (h *taskHandler) Create(c echo.Context) error {
	var taskRequest domain.TaskRequest
	if err := c.Bind(&taskRequest); err != nil {
		return err
	}

	user, err := auth.GetUserFromContext(c)
	if err != nil {
		return err
	}

	_, err = h.service.Create(c.Request().Context(), user.ID, &taskRequest)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Refresh", "true")
	return c.NoContent(201)
}

func (h *taskHandler) GetAllTasks(c echo.Context) error {
	tasks, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	taskListing := task_views.List(task_views.ListProps{Tasks: tasks})
	return api.Render(c, 200, taskListing)
}

func (h *taskHandler) GetForm(c echo.Context) error {
	taskForm := task_views.Form(task_views.FormProps{
		IsEdit: false,
	})
	return api.Render(c, 200, taskForm)
}

type TaskIDParams struct {
	TaskID       int64 `param:"id"`
	ParentTaskID int64 `param:"parent_id"`
}

func (h *taskHandler) GetEditForm(c echo.Context) error {
	var params TaskIDParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	task, err := h.service.GetByID(c.Request().Context(), params.TaskID)
	if err != nil {
		return err
	}

	taskForm := task_views.Form(task_views.FormProps{
		IsEdit: true,
		Task:   task,
	})
	return api.Render(c, 200, taskForm)
}

func (h *taskHandler) Update(c echo.Context) error {
	var params TaskIDParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	var taskRequest domain.TaskRequest
	if err := c.Bind(&taskRequest); err != nil {
		return err
	}

	_, err := h.service.Update(c.Request().Context(), params.TaskID, &taskRequest)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Refresh", "true")
	return c.NoContent(200)
}

func (h *taskHandler) Delete(c echo.Context) error {
	var params TaskIDParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), params.TaskID); err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Refresh", "true")
	return c.NoContent(204)
}

type TaskSelectIDParams struct {
	TaskID       int64 `query:"task_id"`
	ParentTaskID int64 `query:"parent_task_id"`
	ExcludedID   int64 `query:"excluded_id"`
}

func (h *taskHandler) GetSelect(c echo.Context) error {
	var params TaskSelectIDParams
	if err := c.Bind(&params); err != nil {
		return err
	}

	var taskID int64
	if params.TaskID != 0 {
		taskID = params.TaskID
	} else if params.ParentTaskID != 0 {
		taskID = params.ParentTaskID
	}

	tasks, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	taskSelect := task_views.Select(task_views.SelectProps{
		Tasks:          tasks,
		SelectedTaskID: taskID,
		ExcludedID:     params.ExcludedID,
	})
	return api.Render(c, 200, taskSelect)
}
