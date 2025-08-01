package hub

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/ngxxx307/sandbox_vr_wordle/service"
)

// GameSession manages a single two-player wordle game.
type GameSession struct {
	players     [2]*service.MultiplayerWordle
	hub         *Hub
	broadcast   chan *service.ClientMessage
	turn        int
	chances     [2]int
	answer      string
	lookupSet   map[rune]struct{}
	playerIDMap map[*uuid.UUID]int
}

func NewGameSession(hub *Hub, p1, p2 *service.MultiplayerWordle) *GameSession {
	session := &GameSession{
		players:     [2]*service.MultiplayerWordle{p1, p2},
		hub:         hub,
		broadcast:   make(chan *service.ClientMessage),
		turn:        0,
		chances:     [2]int{hub.config.WordleMaxChances, hub.config.WordleMaxChances},
		answer:      hub.Answer, // For simplicity, all games have the same answer for now
		lookupSet:   hub.lookupSet,
		playerIDMap: make(map[*uuid.UUID]int),
	}

	// Map client UUIDs to player numbers (0 or 1)
	session.playerIDMap[&p1.ClientID] = 0
	session.playerIDMap[&p2.ClientID] = 1

	// Crucially, tell each player's service which channel to send messages to.
	p1.ServerChannel = session.broadcast
	p2.ServerChannel = session.broadcast

	return session
}

// Run starts the game loop for the session.
func (s *GameSession) Run() {
	rules := "Welcome to Multiplayer Wordle! A game has been found.\n" +
		"It is a turn-based competitive game! You have 6 tries.\n" +
		"Player 1 starts."
	s.sendToPlayers(&service.ServerMessage{Finished: false, Message: rules})

	for clientMessage := range s.broadcast {
		playerNum, ok := s.playerIDMap[clientMessage.ClientID]
		if !ok {
			continue // Message from a client not in this session.
		}

		if s.turn != playerNum {
			msg := fmt.Sprintf("Server: It is not your turn! It's player %d's turn.", s.turn+1)
			s.players[playerNum].SendChannel <- &service.ServerMessage{
				Message:  msg,
				Finished: false,
			}
			continue
		}

		if len(clientMessage.Message) != 5 {
			msg := fmt.Sprintf("Server: Invalid word length! %s", clientMessage.Message)
			s.players[playerNum].SendChannel <- &service.ServerMessage{
				Message:  msg,
				Finished: false,
			}
			continue
		}

		resp, finished := s.guess(clientMessage.Message, playerNum)
		s.sendToPlayers(&service.ServerMessage{
			Message:  fmt.Sprintf("Player %d guessed: %s -> %s", playerNum+1, clientMessage.Message, resp),
			Finished: finished,
		})

		if finished {
			// Close the broadcast channel to end the loop and the goroutine.
			close(s.broadcast)
			// TODO: Unregister session from hub
			return
		}

		// Switch turns
		s.turn = (s.turn + 1) % 2
		s.sendToPlayers(&service.ServerMessage{
			Message:  fmt.Sprintf("It is now Player %d's turn.", s.turn+1),
			Finished: false,
		})
	}
}

// sendToPlayers broadcasts a message to both players in the session.
func (s *GameSession) sendToPlayers(msg *service.ServerMessage) {
	for _, p := range s.players {
		// Use a select to avoid blocking if a client's channel is full or closed.
		select {
		case p.SendChannel <- msg:
		default:
			// If we can't send, the client might be disconnected.
			// The read/write pump for that client will handle the cleanup.
		}
	}
}

func (s *GameSession) guess(guess string, playerId int) (result string, finished bool) {
	guess = strings.ToUpper(guess)
	if s.chances[playerId] <= 0 {
		return "GAMEOVER", true
	}

	if len(guess) != len(s.answer) {
		return "Invalid Word length!", false
	}

	s.chances[playerId]--

	resultRunes := make([]rune, len(s.answer))
	correct := true

	for i, r := range guess {
		letter := rune(s.answer[i])
		if letter == r {
			resultRunes[i] = 'O'
			continue
		}
		correct = false
		if _, present := s.lookupSet[r]; present {
			resultRunes[i] = '?'
		} else {
			resultRunes[i] = '_'
		}
	}

	if correct {
		return fmt.Sprintf("Bingo! Player %d wins! The answer is %s", playerId+1, s.answer), true
	}

	if s.chances[playerId] <= 0 {
		gameoverMsg := fmt.Sprintf("%s \nGame over for Player %d! The correct answer is %s", string(resultRunes), playerId+1, s.answer)
		// This doesn't end the game for the other player yet.
		// For now, we'll consider the game over for both.
		return gameoverMsg, true
	}

	return string(resultRunes), false
}
