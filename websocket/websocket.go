package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
)

type WebSocket struct {
	pingInterval   time.Duration
	pingTimeout    time.Duration
	pongWait       time.Duration
	maxMessageSize int64
}

func NewWebsocket(c *config.Config) *WebSocket {
	return &WebSocket{
		pingInterval:   time.Duration(c.PingIntervalSec) * time.Second,
		pingTimeout:    time.Duration(c.PingTimeoutSec) * time.Second,
		pongWait:       time.Duration(c.PongWaitSec) * time.Second,
		maxMessageSize: c.MaxMessageSize,
	}
}

// This is likely in your handler/router file, not in the websocket/conn.go file.

func (ws *WebSocket) WebSocketHandler(c echo.Context) error {
	// 1. Create the connection using your new constructor.
	//    This sets up the pong handler and the initial read deadline.
	conn, err := NewConnection(c, ws.pingInterval, ws.pingTimeout, ws.pongWait, ws.maxMessageSize)
	if err != nil {
		log.Println("Failed to upgrade connection in WebSocketHandler:", err)
		return err
	}
	defer conn.Close()

	log.Println("WebSocket Connection Established")

	// 2. === THIS IS THE CRUCIAL, MISSING PIECE ===
	//    Start a goroutine to send pings to the client periodically.
	go func() {
		// Create a ticker that fires every `pingInterval`.
		ticker := time.NewTicker(conn.PingPeriod)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				// When the ticker fires, call the PingPong method on your Conn struct.
				if err := conn.PingPong(); err != nil {
					log.Println("Failed to send ping, closing connection:", err)
					// If ping fails, the connection is likely broken, so we can exit the goroutine.
					// The main read loop will also likely fail and handle the full closure.
					return
				}
			// This case handles the server shutting down or the connection closing.
			case <-c.Request().Context().Done():
				log.Println("Ping goroutine stopping because context is done.")
				return
			}
		}
	}()
	// =======================================================

	// 3. This is the main read loop. It blocks until a message arrives or an error occurs.
	//    The read deadline is now being successfully extended by the pong handler,
	//    which is triggered by the pings sent from the goroutine above.
	for {
		mt, rawMessage, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket connection closed abnormally: %v", err)
			} else {
				log.Println("WebSocket connection closed normally")
			}
			break // Exit the loop on any error.
		}

		log.Printf("Received message: %s", rawMessage)

		// Echo the message back.
		if err := conn.WriteMessage(mt, rawMessage); err != nil {
			log.Println("write error:", err)
			break
		}
	}

	return nil
}
