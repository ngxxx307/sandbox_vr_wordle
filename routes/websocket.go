package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	"github.com/ngxxx307/sandbox_vr_wordle/controller"
)

func SetupWebSocketRoute(e *echo.Echo, svc *controller.WebSocket) {
	webSocketRoute := e.Group("/ws")
	webSocketRoute.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get the config from the DI container
			// This is a bit of a hack, but it's the easiest way to get the config to the handler
			// without changing the method signature.
			cfg, _ := c.Get("config").(*config.Config)
			c.Set("config", cfg)
			return next(c)
		}
	})
	webSocketRoute.GET("/", svc.WebSocketHandler)
}
