package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/ngxxx307/sandbox_vr_wordle/controller"
)

func SetupWebSocketRoute(e *echo.Echo, svc *controller.WebSocket) {
	webSocketRoute := e.Group("/ws")
	webSocketRoute.GET("/", svc.WebSocketHandler)
}
