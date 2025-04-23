package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/auth"
	"github.com/mjmarrazzo/maintenance-app/components/location_views"
	"github.com/mjmarrazzo/maintenance-app/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/api"
	"github.com/mjmarrazzo/maintenance-app/internal/database"
	"github.com/mjmarrazzo/maintenance-app/internal/validation"
	"github.com/mjmarrazzo/maintenance-app/service"
)

type LocationHandler interface {
	api.Handler
	Create(c echo.Context) error
	GetAllLocations(c echo.Context) error
	GetForm(c echo.Context) error
	GetEditForm(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
	GetLocationSelect(c echo.Context) error
}

type locationHandler struct {
	service service.LocationService
}

func (c locationHandler) RegisterRoutes(e *echo.Echo) {
	group := e.Group("/locations")
	group.Use(auth.AuthenticatedMiddleware())

	group.POST("", c.Create)
	group.GET("", c.GetAllLocations)
	group.GET("/form", c.GetForm)
	group.GET("/:id/form", c.GetEditForm)
	group.PUT("/:id", c.Update)
	group.DELETE("/:id", c.Delete)
	group.GET("/select", c.GetLocationSelect)
}

func NewLocationHandler(db *database.Client) LocationHandler {
	return &locationHandler{service: service.NewLocationService(db.Pool())}
}

func (h *locationHandler) Create(c echo.Context) error {
	var locationRequest domain.LocationRequest
	if err := validation.BindBody(c, &locationRequest); err != nil {
		return err
	}

	_, err := h.service.Create(c.Request().Context(), &locationRequest)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Refresh", "true")
	return c.NoContent(201)
}

func (h *locationHandler) GetAllLocations(c echo.Context) error {
	locations, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	locationListing := location_views.List(location_views.ListProps{Locations: locations})
	return api.Render(c, 200, locationListing)
}

func (h *locationHandler) GetForm(c echo.Context) error {
	locations, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	locationForm := location_views.Form(location_views.FormProps{
		IsEdit:       false,
		Location:     nil,
		AllLocations: locations,
	})
	return api.Render(c, 200, locationForm)
}

type LocationIDParam struct {
	ID int64 `param:"location_id" query:"location_id"`
}

func (h *locationHandler) GetEditForm(c echo.Context) error {
	var params LocationIDParam
	if err := c.Bind(&params); err != nil {
		return err
	}

	location, err := h.service.GetByID(c.Request().Context(), params.ID)
	if err != nil {
		return err
	}

	locations, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	locationForm := location_views.Form(location_views.FormProps{
		IsEdit:       true,
		Location:     location,
		AllLocations: locations,
	})
	return api.Render(c, 200, locationForm)
}

func (h *locationHandler) Update(c echo.Context) error {
	var params LocationIDParam
	if err := c.Bind(&params); err != nil {
		return err
	}

	var locationRequest domain.LocationRequest
	if err := c.Bind(&locationRequest); err != nil {
		return err
	}

	_, err := h.service.Update(c.Request().Context(), params.ID, &locationRequest)
	if err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Refresh", "true")
	return c.NoContent(204)
}

func (h *locationHandler) Delete(c echo.Context) error {
	var params LocationIDParam
	if err := c.Bind(&params); err != nil {
		return err
	}

	if err := h.service.Delete(c.Request().Context(), params.ID); err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Refresh", "true")
	return c.NoContent(204)
}

func (h *locationHandler) GetLocationSelect(c echo.Context) error {
	var params LocationIDParam
	if err := c.Bind(&params); err != nil {
		return err
	}

	locations, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		return err
	}

	locationSelect := location_views.Select(location_views.SelectProps{
		Locations:          locations,
		SelectedLocationID: params.ID,
	})
	return api.Render(c, 200, locationSelect)
}
