package controller

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	"github.com/ngxxx307/sandbox_vr_wordle/service"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

type WordleController struct {
	config  *config.Config
	handler service.Handler
}

func NewWordleController(cfg *config.Config) *WordleController {
	svc := service.NewWordleGame(cfg)
	return &WordleController{
		config:  cfg,
		handler: svc,
	}
}

func (wc *WordleController) Handle(conn *w.Conn) Controller {
	rules := "Welcome to Wordle!\n" +
		"You have 6 tries to guess the 5-letter word.\n" +
		"- Each guess must be a valid 5-letter word.\n" +
		"- After each guess, the color of the tiles will change to show how close your guess was to the word.\n" +
		"- O: The letter is in the word and in the correct spot.\n" +
		"- ?: The letter is in the word but in the wrong spot.\n" +
		"- _: The letter is not in the word in any spot.\n\n" +
		"Good luck!"
	if err := conn.WriteMessage(websocket.TextMessage, []byte(rules)); err != nil {
		log.Println("write error:", err)
		return nil
	}

	for {
		_, rawMessage, err := conn.ReadMessage()
		if err != nil {
			w.HandleReadError(err)
			return nil
		}

		resp, finished := wc.handler.Read(string(rawMessage))

		if err := conn.WriteMessage(websocket.TextMessage, []byte(resp)); err != nil {
			log.Println("write error:", err)
			return nil
		}

		if finished {
			// Game over, return to the lounge.
			return NewGameLoungeController(wc.config)
		}
	}
}
