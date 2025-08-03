package service

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"strings"
	"sync"

	"github.com/ngxxx307/sandbox_vr_wordle/config"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

type ClientMessage struct {
	Client  *MultiplayerClient
	Payload string
}

type MultiplayerHost struct {
	Clients       []*MultiplayerClient
	ReadAnswer    chan *ClientMessage // client -> host
	SendBroadcast chan string         // host -> clients
	BroadcastWg   sync.WaitGroup

	Finish []chan struct{}

	chances   []int
	Answer    string
	lookupSet map[rune]struct{}
}

type MultiplayerClient struct {
	Conn        *w.Conn
	Read        chan string         // host -> client
	SendAnswer  chan *ClientMessage // client -> host
	GameStarted chan int

	Finish chan struct{}
}

func NewMultiPlayerClient(conn *w.Conn, read chan string) *MultiplayerClient {
	return &MultiplayerClient{Conn: conn, Read: read}
}

func NewMultiPlayerHost(config *config.Config) *MultiplayerHost {
	read := make(chan *ClientMessage)
	send := make(chan string)

	randomIndex := rand.N(uint(len(config.WordleWordList)))
	answer := config.WordleWordList[randomIndex]

	return &MultiplayerHost{
		ReadAnswer:    read,
		SendBroadcast: send,
		Finish:        []chan struct{}{},

		Answer:    answer,
		chances:   []int{config.WordleMaxChances, config.WordleMaxChances},
		lookupSet: PrepareLookupSet(answer),
	}
}

func (h *MultiplayerHost) End() {
	for _, c := range h.Finish {
		c <- struct{}{}
	}
}

func (h *MultiplayerHost) Close() {
	if h.ReadAnswer != nil {
		close(h.ReadAnswer)
	}
	if h.SendBroadcast != nil {
		close(h.SendBroadcast)
	}
}

func (h *MultiplayerHost) Run() {
	defer func() {
		for _, c := range h.Finish {
			if c != nil {
				close(c)
			}
		}
	}()

	turn := 0
	for {
		clientMsg, ok := <-h.ReadAnswer
		if !ok {
			return // Channel closed, exit goroutine
		}
		index := slices.Index(h.Clients, clientMsg.Client)

		if index == -1 {
			println("Error: invalid client")
			continue
		}
		if turn != index {
			h.SendBroadcast <- fmt.Sprintf("Player %d, it is not your turn!", index)
			continue
		}

		resp, finished := h.Guess(index, clientMsg.Payload)

		// Combine messages into one to avoid blocking on the unbuffered channel
		fullMessage := fmt.Sprintf("Player %d Guess: %s.\n%s", index, clientMsg.Payload, resp)
		h.BroadcastWg.Add(1)
		h.SendBroadcast <- fullMessage

		turn = (turn + 1) % 2

		if finished {
			h.BroadcastWg.Wait()
			h.End()
		}
	}
}

func (h *MultiplayerHost) Read(id int, msg string) (resp string, finished bool) {
	resp, finished = h.Guess(id, msg)
	if finished {
		return resp, true
	}
	return resp, false
}

func (h *MultiplayerHost) Guess(id int, guess string) (result string, finished bool) {
	guess = strings.ToUpper(guess)
	// Game is already over
	if h.chances[id] <= 0 {
		return "GAMEOVER", true
	}

	// Ensure guess is the same length as the answer
	if len(guess) != len(h.Answer) {
		return "Invalid Word length!", false
	}

	h.chances[id]--

	resultRunes := make([]rune, len(h.Answer))
	correct := true

	for i, r := range guess {
		letter := rune(h.Answer[i])
		if letter == r { // correct letter
			resultRunes[i] = 'O'
			continue
		}
		correct = false // If answer is correct, will not run this line
		if _, present := h.lookupSet[r]; present {
			resultRunes[i] = '?'
		} else {
			resultRunes[i] = '_'
		}
	}

	if correct {
		return fmt.Sprintf("bingo! the correct answer is %s", h.Answer), true
	}

	// Check if the game is over due to running out of chances
	for i := 0; i < len(h.Clients); i++ {
		if h.chances[i] > 0 {
			return string(resultRunes), finished
		}
	}
	gameoverMsg := fmt.Sprintf("%s \nGame over! The correct answer is %s", string(resultRunes), h.Answer)
	return gameoverMsg, true
}
