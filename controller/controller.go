package controller

import (
	"github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

type Controller interface {
	Handle(conn *websocket.Conn) Controller
}
