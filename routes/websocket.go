package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

func SetupWebSocketRoute(e *echo.Echo, svc *websocket.WebSocket) {
	websocketRoute := e.Group("/data/ws")
	websocketRoute.GET("/", svc.WebsocketHandler)
}
