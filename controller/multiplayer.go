package controller

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	"github.com/ngxxx307/sandbox_vr_wordle/hub"
	"github.com/ngxxx307/sandbox_vr_wordle/service"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

type MultiplayerWordleController struct {
	config  *config.Config
	Handler *service.MultiplayerWordle
	hub     *hub.Hub
}

func NewMultiplayerWordleController(cfg *config.Config, hub *hub.Hub) *MultiplayerWordleController {
	svc := service.NewMultiplayerWordle(cfg)
	return &MultiplayerWordleController{
		config:  cfg,
		Handler: svc,
		hub:     hub,
	}
}

func (wc *MultiplayerWordleController) Handle(conn *w.Conn) Controller {
	defer conn.Close()

	handler := service.NewMultiplayerWordle(wc.config)
	wc.hub.Enqueue(handler)

	go func() {
		for {
			_, rawMessage, err := conn.ReadMessage()
			if err != nil {
				w.HandleReadError(err)
				close(handler.SendChannel)
				return
			}
			handler.Read(string(rawMessage))
		}
	}()

	for serverMessage := range handler.SendChannel {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(serverMessage.Message)); err != nil {
			log.Println("write error:", err)
			return nil
		}

		if serverMessage.Finished {
			return NewGameLoungeController(wc.config, wc.hub)
		}
	}

	return NewGameLoungeController(wc.config, wc.hub)
}
