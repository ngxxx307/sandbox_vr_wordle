package controller

import (
	"github.com/gorilla/websocket"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

type MultiplayerWordleController struct {
	config *config.Config
}

func NewMultiplayerWordleController(cfg *config.Config) *MultiplayerWordleController {
	return &MultiplayerWordleController{
		config: cfg,
	}
}

func (wc *MultiplayerWordleController) Handle(conn *w.Conn) Controller {
	// Don't start pumps here since they're already started in WebSocketHandler

	for {
		select {
		case msg, ok := <-conn.ReadChannel:
			if !ok {
				// Channel is closed, exit gracefully
				return NewGameLoungeController(wc.config)
			}
			println("msg:", msg.Msg)
			conn.WriteChannel <- &w.WebSocketMessage{Msg: msg.Msg, MessageType: websocket.TextMessage}
			if msg.Msg == "quit" {
				conn.WriteChannel <- &w.WebSocketMessage{Msg: "OKOK", MessageType: websocket.TextMessage}
				return NewGameLoungeController(wc.config)
			}
		}
	}
}
