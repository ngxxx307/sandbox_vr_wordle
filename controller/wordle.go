package controller

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/ngxxx307/sandbox_vr_wordle/service"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

type WordleController struct {
	ctx     *GameContext
	handler *service.Wordle
}

func NewWordleController(ctx *GameContext) *WordleController {
	svc := service.NewWordleGame(ctx.Config)
	return &WordleController{
		handler: svc,
		ctx:     ctx,
	}
}

func (wc *WordleController) Handle(conn *w.Conn) Controller {
	rules := fmt.Sprintf("Welcome to Wordle!\n"+
		"You have %d tries to guess the 5-letter word.\n"+
		"- O: The letter is in the word and in the correct spot.\n"+
		"- ?: The letter is in the word but in the wrong spot.\n"+
		"- _: The letter is not in the word in any spot.\n\n"+
		"Good luck!", wc.ctx.Config.WordleMaxChances)
	conn.WriteChannel <- &w.WebSocketMessage{Msg: rules, MessageType: websocket.TextMessage}

	for {
		msg, ok := <-conn.ReadChannel
		if !ok {
			// Channel is closed, exit gracefully
			return nil
		}

		resp, finished := wc.handler.Read(msg.Msg)
		conn.WriteChannel <- &w.WebSocketMessage{Msg: resp, MessageType: websocket.TextMessage}

		// if finished, relinquish control back to lounge controller
		if finished {
			return NewGameLoungeController(wc.ctx)
		}
	}
}
