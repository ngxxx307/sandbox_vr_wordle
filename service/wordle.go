package service

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/ngxxx307/sandbox_vr_wordle/config"
)

type Wordle struct {
	chances   int
	Answer    string
	lookupSet map[rune]struct{}
	config    *config.Config
}

func PrepareLookupSet(s string) map[rune]struct{} {
	lookupSet := make(map[rune]struct{}, len(s))
	for _, r := range s {
		lookupSet[r] = struct{}{}
	}
	return lookupSet
}

func NewWordleGame(config *config.Config) *Wordle {
	randomIndex := rand.N(uint(len(config.WordleWordList)))
	answer := config.WordleWordList[randomIndex]

	return &Wordle{
		Answer:    answer,
		chances:   config.WordleMaxChances,
		lookupSet: PrepareLookupSet(answer),
		config:    config,
	}
}

func (w *Wordle) Read(msg string) (resp string, finished bool) {
	resp, finished = w.Guess(msg)
	if finished {
		return resp, true
	}
	return resp, false
}

func (w *Wordle) Guess(word string) (result string, finished bool) {
	word = strings.ToUpper(word)
	// Game is already over
	if w.chances <= 0 {
		return "GAMEOVER", true
	}

	// Ensure guess is the same length as the answer
	if len(word) != len(w.Answer) {
		return "Invalid Word length!", false
	}

	w.chances--

	resultRunes := make([]rune, len(w.Answer))
	correct := true

	for i, r := range word {
		letter := rune(w.Answer[i])
		if letter == r { // correct letter
			resultRunes[i] = 'O'
			continue
		}
		correct = false // If answer is correct, will not run this line
		if _, present := w.lookupSet[r]; present {
			resultRunes[i] = '?'
		} else {
			resultRunes[i] = '_'
		}
	}

	if correct {
		return fmt.Sprintf("bingo! the correct answer is %s", w.Answer), true
	}

	// Check if the game is over due to running out of chances
	if w.chances <= 0 {
		finished = true
	}

	return string(resultRunes), finished
}
