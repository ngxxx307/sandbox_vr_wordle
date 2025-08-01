package service

import (
	"fmt"
	"math/rand/v2"

	"github.com/ngxxx307/sandbox_vr_wordle/config"
)

type Wordle struct {
	chances   int
	answer    string
	lookupSet map[rune]struct{}
	config    *config.Config
}

func NewWordleGame(config *config.Config) *Wordle {
	randomIndex := rand.N(uint(len(config.WordleWordList)))
	answer := config.WordleWordList[randomIndex]

	return &Wordle{
		answer:    answer,
		chances:   config.WordleMaxChances,
		lookupSet: PrepareLookupSet(answer),
		config:    config,
	}
}

func (w *Wordle) Read(msg string) (resp string, next Handler) {
	resp, finished := w.Guess(msg)
	if finished {
		// When the game is over, we need to transition back to the default handler.
		// This assumes you have a way to create a new DefaultHandler.
		// If NewDefaultHandler requires dependencies, this part will need adjustment.
		return resp, NewDefaultHandler(w.config)
	}
	return resp, w
}

func (w *Wordle) Guess(word string) (result string, finished bool) {
	// Game is already over
	if w.chances <= 0 {
		return "GAME_OVER", true
	}

	// Ensure guess is the same length as the answer
	if len(word) != len(w.answer) {
		return "INVALID_LENGTH", false
	}

	w.chances--

	// Check for a perfect match (win condition)
	if word == w.answer {
		return fmt.Sprintf("bingo! the correct answer is %s", w.answer), true
	}

	resultRunes := make([]rune, len(w.answer))
	answerFreq := make(map[rune]int)
	for _, r := range w.answer {
		answerFreq[r]++
	}

	// First pass: Find correct letters (Green)
	for i, r := range word {
		if rune(w.answer[i]) == r {
			resultRunes[i] = 'O'
			answerFreq[r]--
		}
	}

	// Second pass: Find present (Yellow) and absent (Gray) letters
	for i, r := range word {
		// Skip letters already marked as correct
		if resultRunes[i] == 'O' {
			continue
		}

		// Check if the letter is in the answer and hasn't been fully accounted for
		if count, ok := answerFreq[r]; ok && count > 0 {
			resultRunes[i] = '?'
			answerFreq[r]--
		} else {
			resultRunes[i] = '_' // Absent
		}
	}

	// Check if the game is over due to running out of chances
	if w.chances <= 0 {
		finished = true
	}

	return string(resultRunes), finished
}
