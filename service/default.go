package service

import (
	"fmt"

	"github.com/ngxxx307/sandbox_vr_wordle/config"
)

type DefaultHandler struct {
	config *config.Config
}

func NewDefaultHandler(cfg *config.Config) Handler {
	return &DefaultHandler{config: cfg}
}

func (h *DefaultHandler) Read(msg string) (resp string, next Handler) {
	switch msg {
	case "Wordle":
		return "Wordle Game start!", NewWordleGame(h.config)
	case "Cheated Wordle":
		return "Wordle Game start!", h // Replace with actual Wordle handler
	case "Multiplayer Wordle":
		return "Wordle Game start!", h // Replace with actual Wordle handler
	}
	return fmt.Sprintf("Error game type: %s", msg), h
}
