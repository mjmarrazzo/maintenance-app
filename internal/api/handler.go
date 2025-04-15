package api

import "github.com/labstack/echo/v4"

type Handler interface {
	RegisterRoutes(g *echo.Echo)
}
