package controller

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/ngxxx307/sandbox_vr_wordle/service"
	w "github.com/ngxxx307/sandbox_vr_wordle/websocket"
)

type MultiplayerWordleController struct {
	ctx *GameContext
}

func NewMultiplayerWordleController(ctx *GameContext) *MultiplayerWordleController {
	return &MultiplayerWordleController{ctx: ctx}
}

func (wc *MultiplayerWordleController) Handle(conn *w.Conn) Controller {
	client := &service.MultiplayerClient{
		Conn:        conn,
		Read:        make(chan string),
		SendAnswer:  make(chan *service.ClientMessage),
		GameStarted: make(chan int),
		Finish:      make(chan struct{}),
	}

	rules := fmt.Sprintf(`Welcome to Multiplayer Wordle!

Rules:
- Each player has to guess a 5-letter word.
- You have %d chances to guess the word.
- 'O': Correct letter, correct position.
- '?': Correct letter, wrong position.
- '_': Incorrect letter.
- Players take turns to guess.

Waiting for another player to join...`, wc.ctx.Config.WordleMaxChances)

	conn.WriteChannel <- &w.WebSocketMessage{Msg: rules, MessageType: websocket.TextMessage}

	wc.ctx.MatchMaker.AddPlayer(client)
	defer wc.ctx.MatchMaker.RemovePlayer(client)

	var wg sync.WaitGroup
	done := make(chan struct{})
	var once sync.Once

	// Goroutine to wait for the finish signal and close the done channel
	go func() {
		<-client.Finish
		once.Do(func() { close(done) })
	}()

	wg.Add(2)

	// Goroutine 1: Read from host and write to websocket
	go func() {
		defer wg.Done()
		for {
			select {
			case msg, ok := <-client.Read:
				if !ok {
					return // Channel closed
				}
				conn.WriteChannel <- &w.WebSocketMessage{Msg: msg, MessageType: websocket.TextMessage}
			case <-done:
				return
			}
		}
	}()

	// Goroutine 2: Read from websocket and write to host
	go func() {
		defer wg.Done()
		// Wait for the game to start before reading input
		select {
		case playerIndex := <-client.GameStarted:
			conn.WriteChannel <- &w.WebSocketMessage{Msg: fmt.Sprintf("Game Started! You are player %d", playerIndex), MessageType: websocket.TextMessage}
		case <-done:
			return
		}

		for {
			select {
			case msg, ok := <-conn.ReadChannel:
				if !ok {
					once.Do(func() { close(done) }) // Websocket closed, signal others to exit
					return
				}
				client.SendAnswer <- &service.ClientMessage{Client: client, Payload: msg.Msg}
			case <-done:
				return
			}
		}
	}()

	wg.Wait() // Wait for both goroutines to finish
	return NewGameLoungeController(wc.ctx)
}
