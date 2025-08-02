package main

import "github.com/ngxxx307/sandbox_vr_wordle/service"

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	Host service.MultiplayerHost
	// Registered clients.
	clients map[*service.MultiplayerClinet]bool

	// Inbound messages from the clients.
	broadcast chan string

	// Register requests from the clients.
	register chan *service.MultiplayerClinet

	// Unregister requests from clients.
	unregister chan *service.MultiplayerClinet
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan string),
		register:   make(chan *service.MultiplayerClinet),
		unregister: make(chan *service.MultiplayerClinet),
		clients:    make(map[*service.MultiplayerClinet]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}
