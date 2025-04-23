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
			authCtx, err := GetAuthContext(c)
			if err != nil {
				return handleUnauthorized(c, err)
			}

			now := time.Now().Unix()
			if now > authCtx.ExpiresAt {
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

func GetAuthContext(c echo.Context) (*AuthContext, error) {
	user, err := getUserFromSession(c)
	if err != nil {
		return nil, err
	}
	expiration, err := getExpirationTime(c)
	if err != nil {
		return nil, err
	}
	authContext := &AuthContext{
		User:      user,
		ExpiresAt: expiration,
	}
	return authContext, nil
}
