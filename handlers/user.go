package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/auth"
	"github.com/mjmarrazzo/maintenance-app/components/auth_views"
	"github.com/mjmarrazzo/maintenance-app/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/api"
	"github.com/mjmarrazzo/maintenance-app/internal/database"
	"github.com/mjmarrazzo/maintenance-app/internal/validation"
	"github.com/mjmarrazzo/maintenance-app/service"
)

type UserHandler interface {
	api.Handler
	Login(c echo.Context) error
	Logout(c echo.Context) error
	GetRegisterForm(c echo.Context) error
	Register(c echo.Context) error
}

type userHandler struct {
	service service.UserService
}

func (h *userHandler) RegisterRoutes(e *echo.Echo) {
	e.POST("/login", h.Login)
	e.GET("/logout", h.Logout)
	e.GET("/register", h.GetRegisterForm)
	e.POST("/register", h.Register)
}

func NewAuthHandler(db *database.Client) UserHandler {
	return &userHandler{
		service: service.NewUserService(db.Pool()),
	}
}

type LoginParams struct {
	Email       string `form:"email" validate:"required,email"`
	Password    string `form:"password" validate:"required"`
	OriginalUrl string `form:"original_url"`
}

func (h *userHandler) Login(c echo.Context) error {
	loginParams := &LoginParams{}
	if err := validation.BindBody(c, loginParams); err != nil {
		return err
	}

	user, err := h.service.Authenticate(c.Request().Context(), loginParams.Email, loginParams.Password)
	if err != nil {
		return err
	}

	if err := auth.SaveUserToSession(c, user); err != nil {
		return err
	}

	redirectUrl := loginParams.OriginalUrl
	if redirectUrl == "" {
		redirectUrl = "/home"
	}

	c.Response().Header().Set("Hx-Redirect", redirectUrl)
	return nil
}

func (h *userHandler) Logout(c echo.Context) error {
	if err := auth.ClearSession(c); err != nil {
		return err
	}
	c.Response().Header().Set("Hx-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

type RegistrationParams struct {
	Email string `query:"email" `
}

func (h *userHandler) GetRegisterForm(c echo.Context) error {
	params := &RegistrationParams{}
	if err := validation.BindPathParams(c, params); err != nil {
		return err
	}

	registerForm := auth_views.Register(auth_views.RegisterProps{
		Email: params.Email,
	})
	return api.Render(c, http.StatusOK, registerForm)
}

func (h *userHandler) Register(c echo.Context) error {
	userRequest := &domain.UserRequest{}
	if err := validation.BindBody(c, userRequest); err != nil {
		return err
	}

	if err := h.service.Create(c.Request().Context(), userRequest); err != nil {
		return err
	}

	c.Response().Header().Set("Hx-Redirect", "/home")

	return nil
}
