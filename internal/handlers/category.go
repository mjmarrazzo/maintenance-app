package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/internal/api"
	"github.com/mjmarrazzo/maintenance-app/internal/auth"
	"github.com/mjmarrazzo/maintenance-app/internal/database"
	"github.com/mjmarrazzo/maintenance-app/internal/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/service"
	"github.com/mjmarrazzo/maintenance-app/internal/views/category_views"
)

type CategoryHandler interface {
	api.Handler
	Create(c echo.Context) error
	GetAllCategories(c echo.Context) error
	GetForm(c echo.Context) error
	GetEditForm(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	GetCategorySelect(c echo.Context) error
}

type categoryHandler struct {
	service service.CategoryService
}

func (c categoryHandler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/categories")
	group.Use(auth.AuthenticatedMiddleware())

	group.POST("", c.Create)
	group.GET("", c.GetAllCategories)
	group.GET("/form", c.GetForm)
	group.GET("/:id/form", c.GetEditForm)
	group.PUT("/:id", c.Update)
	group.DELETE("/:id", c.Delete)
	group.GET("/select", c.GetCategorySelect)
}

func NewCategoryHandler(db *database.Client) CategoryHandler {
	return &categoryHandler{service: service.NewCategoryService(db.Pool())}
}

func (h *categoryHandler) Create(c echo.Context) error {
	var categoryRequest domain.CategoryRequest
	if err := c.Bind(&categoryRequest); err != nil {
		return err
	}

	_, err := h.service.Create(c.Request().Context(), &categoryRequest)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Refresh", "true")
	return c.NoContent(201)
}

func (h *categoryHandler) GetAllCategories(c echo.Context) error {
	categories, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	categoryListing := category_views.List(category_views.ListProps{Categories: categories})
	return api.Render(c, 200, categoryListing)
}

func (h *categoryHandler) GetForm(c echo.Context) error {
	categoryForm := category_views.Form(category_views.FormProps{
		IsEdit:   false,
		Category: nil,
	})
	return api.Render(c, 200, categoryForm)
}

type CategoryIDParam struct {
	ID int64 `param:"id" query:"category_id"`
}

func (h *categoryHandler) GetEditForm(c echo.Context) error {
	var params CategoryIDParam
	if err := c.Bind(&params); err != nil {
		return err
	}

	category, err := h.service.GetByID(c.Request().Context(), params.ID)
	if err != nil {
		return err
	}

	categoryForm := category_views.Form(category_views.FormProps{
		IsEdit:   true,
		Category: category,
	})
	return api.Render(c, 200, categoryForm)
}

func (h *categoryHandler) Update(c echo.Context) error {
	var params CategoryIDParam
	if err := c.Bind(&params); err != nil {
		return err
	}

	var category domain.CategoryRequest
	if err := c.Bind(&category); err != nil {
		return err
	}

	_, err := h.service.Update(c.Request().Context(), params.ID, &category)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Refresh", "true")
	return c.NoContent(204)
}

func (h *categoryHandler) Delete(c echo.Context) error {
	var params CategoryIDParam
	if err := c.Bind(&params); err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), params.ID); err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Refresh", "true")
	return c.NoContent(204)
}

func (h *categoryHandler) GetCategorySelect(c echo.Context) error {
	var params CategoryIDParam
	if err := c.Bind(&params); err != nil {
		return err
	}

	categories, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	categorySelect := category_views.Select(category_views.SelectProps{
		Categories:         categories,
		SelectedCategoryID: params.ID,
	})
	return api.Render(c, 200, categorySelect)
}
