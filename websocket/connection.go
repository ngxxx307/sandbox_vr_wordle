package websocket

import (
	"compress/flate"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
)

type Conn struct {
	conn           *websocket.Conn
	lock           sync.Mutex
	PingPeriod     time.Duration
	PingTimeout    time.Duration
	PongWait       time.Duration
	MaxMessageSize int64
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

	return &Conn{
		conn:           conn,
		PingPeriod:     pingPeriod,
		PingTimeout:    pingTimeout,
		PongWait:       pongWait,
		MaxMessageSize: maxMessageSize,
	}, nil
}

func (c *Conn) Close() error {
	return c.conn.Close()
}

func (c *Conn) WriteJSON(v interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteJSON(v)
}

func (c *Conn) WriteMessage(data string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(websocket.TextMessage, []byte(data))
}

func (c *Conn) PingPong() error {
	deadline := time.Now().Add(c.PingTimeout)
	if err := c.conn.WriteControl(websocket.PingMessage, []byte{}, deadline); err != nil {
		return err
	}
	return nil
}

func (c *Conn) ReadMessage() (string, error) {
	_, rawMessage, err := c.conn.ReadMessage()
	return string(rawMessage), err
}
