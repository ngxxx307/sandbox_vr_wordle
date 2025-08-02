package websocket

// TODO: Revamp the websocket package to be channel-based rather than procedural-based

import (
	"bytes"
	"compress/flate"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	writeWait = 10 * time.Second
)

type WebSocketMessage struct {
	Msg         string
	MessageType int
}

type Conn struct {
	conn           *websocket.Conn
	lock           sync.Mutex
	closed         bool
	PingPeriod     time.Duration
	PingTimeout    time.Duration
	PongWait       time.Duration
	MaxMessageSize int64

	ReadChannel  chan *WebSocketMessage
	WriteChannel chan *WebSocketMessage
}

func NewConnection(c echo.Context, pingPeriod time.Duration, pingTimeout time.Duration, pongWait time.Duration, maxMessageSize int64) (*Conn, error) {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return nil, err
	}
	conn.EnableWriteCompression(true)
	conn.SetCompressionLevel(flate.BestCompression)
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		println("pong")
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	readChannel := make(chan *WebSocketMessage)
	writeChannel := make(chan *WebSocketMessage)

	return &Conn{
		conn:           conn,
		PingPeriod:     pingPeriod,
		PingTimeout:    pingTimeout,
		PongWait:       pongWait,
		MaxMessageSize: maxMessageSize,

		ReadChannel:  readChannel,
		WriteChannel: writeChannel,
	}, nil
}

func (c *Conn) Close() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.closed {
		return nil // Already closed
	}
	c.closed = true

	// Close channels to signal pumps to stop
	if c.ReadChannel != nil {
		close(c.ReadChannel)
	}
	if c.WriteChannel != nil {
		close(c.WriteChannel)
	}

	return c.conn.Close()
}

func (c *Conn) PingPong() error {
	c.lock.Lock()
	defer c.lock.Unlock()
	deadline := time.Now().Add(c.PingTimeout)
	if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, deadline); err != nil {
		return err
	}
	return nil
}

func (c *Conn) ReadPump() {
	defer func() {
		// Only close if not already closed
		c.lock.Lock()
		if !c.closed {
			c.conn.Close()
		}
		c.lock.Unlock()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("ReadPump error: %v", err)
			}
			return
		}

		message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))

		// Try to send message, but return if channel is closed
		select {
		case c.ReadChannel <- &WebSocketMessage{Msg: string(message), MessageType: websocket.TextMessage}:
		default:
			// Channel is likely closed or full, exit gracefully
			return
		}
	}
}

func (c *Conn) WritePump() {
	defer func() {
		// Only close if not already closed
		c.lock.Lock()
		if !c.closed {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			c.conn.Close()
		}
		c.lock.Unlock()
	}()

	for message := range c.WriteChannel {
		c.lock.Lock()

		// Check if connection is closed
		if c.closed {
			c.lock.Unlock()
			return
		}

		c.conn.SetWriteDeadline(time.Now().Add(writeWait))

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			c.lock.Unlock()
			return
		}
		w.Write([]byte(message.Msg))

		// Add queued chat messages to the current websocket message.
		n := len(c.WriteChannel)
		for i := 0; i < n; i++ {
			w.Write([]byte("\n"))
			message := <-c.WriteChannel
			w.Write([]byte(message.Msg))
		}

		if err := w.Close(); err != nil {
			c.lock.Unlock()
			return
		}
		c.lock.Unlock()
	}
}
