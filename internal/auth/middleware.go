package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/internal/api"
	"github.com/mjmarrazzo/maintenance-app/internal/domain"
	"github.com/mjmarrazzo/maintenance-app/internal/views/auth_views"
)

type AuthContext struct {
	User      *domain.User
	ExpiresAt int64
}

func AuthenticatedMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			_, err := getUserFromSession(c)
			if err != nil {
				return handleUnauthorized(c, err)
			}

			expiration, err := getExpirationTime(c)
			if err != nil {
				return handleUnauthorized(c, err)
			}

			now := time.Now().Unix()
			if now > expiration {
				return handleUnauthorized(c, err)
			}

			return next(c)
		}
	}
}

func handleUnauthorized(c echo.Context, err error) error {
	log.Printf("Unauthorized access: %v\n", err)
	login := auth_views.Login()
	return api.Render(c, http.StatusUnauthorized, login)
}

func GetUserFromContext(c echo.Context) (*domain.User, error) {
	user, err := getUserFromSession(c)
	if err != nil {
		return nil, err
	}
	return user, nil
}
