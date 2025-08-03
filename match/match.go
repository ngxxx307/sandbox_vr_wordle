package match

import "github.com/ngxxx307/sandbox_vr_wordle/service"

type match struct {
	Host *service.MultiplayerHost
	// Registered clients.
	clients map[*service.MultiplayerClient]bool

	// Inbound messages from the clients.
	broadcast chan string

	// Inbound messages from the clients for the host.
	toHost chan *service.ClientMessage

	// Register requests from the clients.
	Register chan *service.MultiplayerClient

	// Unregister requests from clients.
	Unregister chan *service.MultiplayerClient

	Finish chan struct{}
}

func NewMatch(host *service.MultiplayerHost) *match {
	broadcast := make(chan string, 10)

	host.SendBroadcast = broadcast
	return &match{
		Host:       host,
		broadcast:  broadcast,
		toHost:     make(chan *service.ClientMessage),
		Register:   make(chan *service.MultiplayerClient),
		Unregister: make(chan *service.MultiplayerClient),
		clients:    make(map[*service.MultiplayerClient]bool),
		Finish:     make(chan struct{}),
	}
}

func (h *match) Run() {
	defer func() {
		// Clean up all client connections when the match ends
		for client := range h.clients {
			if client.Read != nil {
				close(client.Read)
			}
			if client.Finish != nil {
				close(client.Finish)
			}
			delete(h.clients, client)
		}
	}()
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				if client.Read != nil {
					close(client.Read)
				}
			}
		case message := <-h.Host.SendBroadcast:
			for client := range h.clients {
				select {
				case client.Read <- message:
				default:
					// Assume client is disconnected, and unregister them
					go func(c *service.MultiplayerClient) {
						h.Unregister <- c
					}(client)
				}
			}
			h.Host.BroadcastWg.Done()
		case <-h.Finish:
			return
		}
	}
}
