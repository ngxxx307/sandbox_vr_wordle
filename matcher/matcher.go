package matcher

// import (
// 	"sync"
// 	"time"

// 	"github.com/ngxxx307/sandbox_vr_wordle/config"
// 	"github.com/ngxxx307/sandbox_vr_wordle/service"
// )

// type Matcher struct {
// 	Queue       []*service.MultiplayerWordle
// 	QueueLock   sync.Mutex
// 	QueueUpdate chan struct{}
// 	stop        chan struct{}
// 	config      *config.Config
// }

// func NewMatcher() *Matcher {
// 	return &Matcher{}
// }

// func (m *Matcher) RunMatchMaking() {
// 	ticker := time.NewTicker(5 * time.Second)
// 	for {
// 		select {
// 		case <-m.QueueUpdate:
// 			m.CheckMatchMaking()
// 		case <-ticker.C:
// 			m.CheckMatchMaking()
// 		case <-m.stop:
// 			return
// 		}
// 	}
// }

// func (m *Matcher) CheckMatchMaking() {
// 	// Loop until n_queue < 2
// 	for len(m.Queue) >= 2 {
// 		m.QueueLock.Lock()
// 		p1 := m.Queue[0]
// 		p2 := m.Queue[1]
// 		m.Queue = m.Queue[2:]
// 		m.QueueLock.Unlock()

// 		// Create and run a new game session for the matched players.
// 		session := NewGameSession(m, p1, p2)
// 		go session.Run()
// 	}
// }
