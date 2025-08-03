package matchMaker

import (
	"fmt"
	"sync"
	"time"

	"github.com/ngxxx307/sandbox_vr_wordle/config"
	match "github.com/ngxxx307/sandbox_vr_wordle/hub"
	"github.com/ngxxx307/sandbox_vr_wordle/service"
)

type MatchMaker struct {
	Queue       []*service.MultiplayerClient
	Lock        sync.Mutex
	QueueUpdate chan struct{}
}

func NewMatcher() *MatchMaker {

	return &MatchMaker{
		Queue:       make([]*service.MultiplayerClient, 0),
		QueueUpdate: make(chan struct{}),
		Lock:        sync.Mutex{},
	}
}

func RunMatchMaking(m *MatchMaker, c *config.Config) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-m.QueueUpdate:
				m.CheckMatchMaking(c)
			case <-ticker.C:
				m.CheckMatchMaking(c)
			}
		}
	}()
}

func (m *MatchMaker) CheckMatchMaking(c *config.Config) {
	m.Lock.Lock()
	defer m.Lock.Unlock()
	for len(m.Queue) >= 2 {

		host := service.NewMultiPlayerHost(c)
		session := match.NewMatch(host)

		clients := make([]*service.MultiplayerClient, 2)
		for i := 0; i < 2; i++ {
			p := m.Queue[0]
			m.Queue = m.Queue[1:]
			clients[i] = p // Assign player to the slice index
			if host.ReadAnswer == nil {
				fmt.Println("Error: host read channel is nil")
			}
			p.SendAnswer = host.ReadAnswer
			host.Finish = append(host.Finish, p.Finish)
		}
		host.Clients = clients
		go host.Run()
		go session.Run()

		for i, c := range clients {
			session.Register <- c
			c.GameStarted <- i
		}
	}
}

func (m *MatchMaker) AddPlayer(p *service.MultiplayerClient) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	m.Queue = append(m.Queue, p)
	select {
	case m.QueueUpdate <- struct{}{}:
	default:
		// Channel is full, skip signal
	}
}

// RemovePlayer removes a player from the matchmaking queue
func (m *MatchMaker) RemovePlayer(player *service.MultiplayerClient) {
	m.Lock.Lock()
	defer m.Lock.Unlock()

	for i, p := range m.Queue {
		if p == player {
			m.Queue = append(m.Queue[:i], m.Queue[i+1:]...)
			break
		}
	}
}

func (m *MatchMaker) GetQueueSize() int {
	return len(m.Queue)
}
