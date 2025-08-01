package controller

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	"github.com/ngxxx307/sandbox_vr_wordle/hub"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

// GameLoungeController is the initial controller that lets users select a game.
type GameLoungeController struct {
	config *config.Config
	hub    *hub.Hub
}

func NewGameLoungeController(cfg *config.Config, hub *hub.Hub) *GameLoungeController {
	return &GameLoungeController{
		config: cfg,
		hub:    hub,
	}
}

// Handle manages the connection while in the game lounge.
func (g *GameLoungeController) Handle(conn *w.Conn) Controller {
	infoMsg := "Welcome! Available games: Wordle, Cheated Host Wordle, Multiplayer Wordle."
	if err := conn.WriteMessage(websocket.TextMessage, []byte(infoMsg)); err != nil {
		log.Println("Failed to send welcome message:", err)
		return nil
	}

	_, rawMessage, err := conn.ReadMessage()
	if err != nil {
		w.HandleReadError(err)
		return nil
	}
	var req = string(rawMessage)
	var resp string
	var nextController Controller

	switch req {
	case "Wordle":
		resp = "Wordle Game start!"
		nextController = NewWordleController(g.config, g.hub)
	case "Cheated Host Wordle":
		resp = "Cheated Host Wordle game start!"
		nextController = NewCheatedHostController(g.config, g.hub)
	case "Multiplayer Wordle":
		resp = "Entering multiplayer queue..."
		nextController = NewMultiplayerWordleController(g.config, g.hub)
	default:
		resp = fmt.Sprintf("Error game type: %s", req)
		nextController = g
	}

	if err := conn.WriteMessage(websocket.TextMessage, []byte(resp)); err != nil {
		log.Println("write error:", err)
		return nil
	}
	return nextController

}
