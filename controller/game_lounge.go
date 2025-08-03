package controller

import (
	"fmt"

	"github.com/gorilla/websocket"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

// GameLoungeController is the initial controller that lets users select a game.
type GameLoungeController struct {
	ctx *GameContext
}

func NewGameLoungeController(ctx *GameContext) *GameLoungeController {
	return &GameLoungeController{
		ctx: ctx,
	}
}

// Handle manages the connection while in the game lounge.
func (c *GameLoungeController) Handle(conn *w.Conn) Controller {
	infoMsg := "Welcome! Available games: Wordle, Cheated Host Wordle, Multiplayer Wordle."
	conn.WriteChannel <- &w.WebSocketMessage{Msg: infoMsg, MessageType: websocket.TextMessage}

	msg, ok := <-conn.ReadChannel
	if !ok {
		// Channel is closed, exit gracefully
		return nil
	}

	var resp string
	var nextController Controller

	switch msg.Msg {
	case "Wordle":
		resp = "Wordle Game start!"
		nextController = NewWordleController(c.ctx)
	case "Cheated Host Wordle":
		resp = "Cheated Host Wordle game start!"
		nextController = NewCheatedHostController(c.ctx)
	case "Multiplayer Wordle":
		resp = "Entering multiplayer queue..."
		nextController = NewMultiplayerWordleController(c.ctx)
	default:
		resp = fmt.Sprintf("Error game type: %s", msg.Msg)
		nextController = c
	}

	conn.WriteChannel <- &w.WebSocketMessage{Msg: resp, MessageType: websocket.TextMessage}
	return nextController

}
