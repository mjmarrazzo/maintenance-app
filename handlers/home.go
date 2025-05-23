package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/components/auth_views"
	"github.com/mjmarrazzo/maintenance-app/components/home_views"
	"github.com/mjmarrazzo/maintenance-app/internal/api"
)

type HomeHandler interface {
	api.Handler
	Index(c echo.Context) error
	Home(c echo.Context) error
}

type homeHandler struct{}

func (h *homeHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/", h.Index)
	e.GET("/home", h.Home)
}

func NewHomeHandler() HomeHandler {
	return &homeHandler{}
}

func (h *homeHandler) Index(c echo.Context) error {
	login := auth_views.Login(auth_views.LoginProps{})
	return api.Render(c, http.StatusOK, login)
}

func (h *homeHandler) Home(c echo.Context) error {
	home := home_views.Home()
	return api.Render(c, http.StatusOK, home)
}
