package api

import (
	"github.com/labstack/echo/v4"
	"github.com/mjmarrazzo/maintenance-app/internal/responses"
)

func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			err := next(e)

			if err == nil {
				return nil
			}

			if ve, ok := responses.IsValidationError(err); ok {
				return e.JSON(400, ve)
			}

			if he, ok := err.(*echo.HTTPError); ok {
				if he.Code == 404 {
					return e.JSON(404, responses.NewNotFoundError("Route not found"))
				}
			}

			e.Logger().Error(err)
			return e.JSON(500, responses.NewInternalServerError("Internal server error"))
		}
	}
}
