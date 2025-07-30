package websocket

import (
	"compress/flate"
	"sync"

	"github.com/gorilla/websocket"
)

type Conn struct {
	conn *websocket.Conn
	lock sync.Mutex
}

func New(conn *websocket.Conn) *Conn {
	conn.EnableWriteCompression(true)
	conn.SetCompressionLevel(flate.BestCompression)
	return &Conn{conn: conn}
}

func (c *Conn) Close() error {
	return c.conn.Close()
}
func (c *Conn) WriteJSON(v interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	return c.conn.WriteJSON(v)
}

func (c *Conn) ReadMessage() (int, []byte, error) {
	return c.conn.ReadMessage()
}
