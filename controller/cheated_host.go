package controller

import (
	"github.com/gorilla/websocket"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	"github.com/ngxxx307/sandbox_vr_wordle/service"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

type CheatedHostController struct {
	config  *config.Config
	handler *service.CheatedHost
}

func NewCheatedHostController(cfg *config.Config) *CheatedHostController {
	svc := service.NewCheatedHostGame(cfg)
	return &CheatedHostController{
		config:  cfg,
		handler: svc,
	}
}

func (wc *CheatedHostController) Handle(conn *w.Conn) Controller {
	rules := "Welcome to Wordle!\n" +
		"You have 6 tries to guess the 5-letter word.\n" +
		"- O: The letter is in the word and in the correct spot.\n" +
		"- ?: The letter is in the word but in the wrong spot.\n" +
		"- _: The letter is not in the word in any spot.\n\n" +
		"Good luck!"
	conn.WriteChannel <- &w.WebSocketMessage{Msg: rules, MessageType: websocket.TextMessage}

	for {
		msg, ok := <-conn.ReadChannel
		if !ok {
			// Channel is closed, exit gracefully
			return NewGameLoungeController(wc.config)
		}

		resp, finished := wc.handler.Read(msg.Msg)
		conn.WriteChannel <- &w.WebSocketMessage{Msg: resp, MessageType: websocket.TextMessage}

		// if finished, relinquish control back to lounge controller
		if finished {
			return NewGameLoungeController(wc.config)
		}
	}
}
