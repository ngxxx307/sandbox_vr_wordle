package controller

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

// GameLoungeController is the initial controller that lets users select a game.
type GameLoungeController struct {
	config *config.Config
}

func NewGameLoungeController(cfg *config.Config) *GameLoungeController {
	return &GameLoungeController{
		config: cfg,
	}
}

// Handle manages the connection while in the game lounge.
func (g *GameLoungeController) Handle(conn *w.Conn) Controller {
<<<<<<< HEAD
	infoMsg := "Welcome! Available games: Wordle, Cheated Host Wordle, Multiplayer Wordle."
=======
	infoMsg := "Welcome! Available games: Wordle, Cheated Wordle, Multiplayer Wordle."
>>>>>>> main
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
		nextController = NewWordleController(g.config)
	case "Cheated Host Wordle":
		resp = "Cheated Host Wordle game start!"
		nextController = NewCheatedHostController(g.config)
	case "Multiplayer Wordle":
		resp = "Not available yet! Coming soon..."
		nextController = g
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
