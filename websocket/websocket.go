package websocket

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
)

type WebSocket struct {
	pingInterval time.Duration
}

func NewWebsocket(c *config.Config) *WebSocket {
	return &WebSocket{}
}

func (ws *WebSocket) WebsocketHandler(c echo.Context) error {

	return nil
}

func (ws *WebSocket) Pingpong(conn *Conn, lastPingtime *time.Time) {
	ticker := time.NewTicker(ws.pingInterval * 2)
	defer ticker.Stop()
	for {
		if lastPingtime != nil && time.Since(*lastPingtime) > ws.pingInterval {
			if err := conn.Close(); err != nil {
				// TODO: handle error
			}
			return
		}
	}
}
