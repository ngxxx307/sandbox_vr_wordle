package service

import (
	"github.com/google/uuid"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
)

type ClientMessage struct {
	ClientID *uuid.UUID
	Message  string
}
type ServerMessage struct {
	Finished bool
	Message  string
}

type MultiplayerWordle struct {
	config        *config.Config
	ClientID      uuid.UUID
	ServerChannel chan *ClientMessage
	SendChannel   chan *ServerMessage // Add a channel for outgoing messages
}

func NewMultiplayerWordle(c *config.Config) *MultiplayerWordle {
	return &MultiplayerWordle{
		config:        c,
		ClientID:      uuid.New(),
		ServerChannel: make(chan *ClientMessage),
		SendChannel:   make(chan *ServerMessage, 1), // Buffered channel
	}
}

// Read is called when a message is received from the client.
// It sends the message to the central hub for processing.
func (w *MultiplayerWordle) Read(msg string) {
	defer func() {
		// This recover is a safeguard. If the channel is closed between the
		// check and the send, this will prevent the server from crashing.
		_ = recover()
	}()

	select {
	case w.ServerChannel <- &ClientMessage{ClientID: &w.ClientID, Message: msg}:
		// Message sent successfully.
	default:
		// The channel is likely closed or full. We can ignore the message.
		// This happens if the game session has already ended.
	}
}
