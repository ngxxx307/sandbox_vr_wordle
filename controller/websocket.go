package controller

import (
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

type WebSocket struct {
	pingInterval   time.Duration
	pingTimeout    time.Duration
	pongWait       time.Duration
	maxMessageSize int64
	config         *config.Config
}

func NewWebsocket(c *config.Config) *WebSocket {
	return &WebSocket{
		pingInterval:   time.Duration(c.PingIntervalSec) * time.Second,
		pingTimeout:    time.Duration(c.PingTimeoutSec) * time.Second,
		pongWait:       time.Duration(c.PongWaitSec) * time.Second,
		maxMessageSize: c.MaxMessageSize,
		config:         c,
	}
}

func (ws *WebSocket) WebSocketHandler(c echo.Context) error {
	conn, err := w.NewConnection(c, ws.pingInterval, ws.pingTimeout, ws.pongWait, ws.maxMessageSize)
	if err != nil {
		log.Println("Failed to upgrade connection in WebSocketHandler:", err)
		return err
	}
	defer conn.Close()

	log.Println("WebSocket Connection Established")

	go func() {
		ticker := time.NewTicker(conn.PingPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := conn.PingPong(); err != nil {
					log.Println("Failed to send ping, closing connection:", err)
					return
				}
			case <-c.Request().Context().Done():
				log.Println("Ping goroutine stopping because context is done.")
				return
			}
		}
	}()

	// Start with the GameLounge controller
	var currentController Controller = NewGameLoungeController(ws.config)

	for currentController != nil {
		nextController := currentController.Handle(conn)
		currentController = nextController
	}

	log.Println("Connection handler loop finished.")
	return nil
}
