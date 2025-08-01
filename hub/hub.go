package hub

import (
	"math/rand/v2"
	"sync"

	"github.com/ngxxx307/sandbox_vr_wordle/config"
	"github.com/ngxxx307/sandbox_vr_wordle/service"
)

// Hub manages the client queue and starts new game sessions.
type Hub struct {
	Queue       []*service.MultiplayerWordle
	QueueLock   sync.Mutex
	QueueUpdate chan struct{}
	stop        chan struct{}
	config      *config.Config

	// These are now part of the GameSession, but the hub can generate them.
	Answer    string
	lookupSet map[rune]struct{}
}

func NewHub(c *config.Config) *Hub {
	// The hub can still be responsible for picking a word for new games.
	randomIndex := rand.N(uint(len(c.WordleWordList)))
	answer := c.WordleWordList[randomIndex]

	return &Hub{
		config:      c,
		Answer:      answer,
		lookupSet:   service.PrepareLookupSet(answer),
		QueueUpdate: make(chan struct{}, 1),
		stop:        make(chan struct{}),
	}
}

// Run starts the hub's main event loop for matchmaking.
func (h *Hub) Run() {
	for {
		select {
		case <-h.QueueUpdate:
			if len(h.Queue) >= 2 {
				h.QueueLock.Lock()
				p1 := h.Queue[0]
				p2 := h.Queue[1]
				h.Queue = h.Queue[2:]
				h.QueueLock.Unlock()

				// Create and run a new game session for the matched players.
				session := NewGameSession(h, p1, p2)
				go session.Run()
			}
		case <-h.stop:
			return
		}
	}
}

// Stop gracefully shuts down the hub.
func (h *Hub) Stop() {
	close(h.stop)
}

// Enqueue adds a client to the matchmaking queue.
func (h *Hub) Enqueue(c *service.MultiplayerWordle) {
	h.QueueLock.Lock()
	h.Queue = append(h.Queue, c)
	h.QueueUpdate <- struct{}{}
	h.QueueLock.Unlock()
}
