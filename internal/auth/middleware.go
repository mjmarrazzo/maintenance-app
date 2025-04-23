package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/internal/api"
	"github.com/mjmarrazzo/maintenance-app/internal/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/views/auth_views"
)

type AuthContext struct {
	User *domain.User
}

func AuthenticatedMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user, err := getUserFromSession(c)
			if err != nil {
				login := auth_views.Login()
				return api.Render(c, http.StatusOK, login)
			}

			if err := SaveUserToSession(c, user); err != nil {
				return err
			}

			return next(c)
		}
	}
}

func GetUserFromContext(c echo.Context) (*domain.User, error) {
	user, err := getUserFromSession(c)
	if err != nil {
		return nil, err
	}
	return user, nil
}
