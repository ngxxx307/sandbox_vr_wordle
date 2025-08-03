package service

import (
	"fmt"
	"math"
	"strings"

	"github.com/ngxxx307/sandbox_vr_wordle/config"
)

type CheatedHost struct {
	chances       int
	WordList      []string
	wordLookupSet map[string]map[rune]struct{} // Example: {"APPLE": {"A", "P", "L", "E"}}
}

func NewCheatedHostGame(config *config.Config) *CheatedHost {
	wordLookupSet := make(map[string]map[rune]struct{}, len(config.WordleWordList))

	for _, word := range config.WordleWordList {
		wordLookupSet[word] = PrepareLookupSet(word)
	}

	return &CheatedHost{
		WordList:      config.WordleWordList,
		wordLookupSet: wordLookupSet,
		chances:       config.WordleMaxChances,
	}
}

func (w *CheatedHost) Read(msg string) (resp string, finished bool) {
	resp, finished = w.Guess(msg)
	if finished {
		return resp, true
	}
	return resp, false
}

func (w *CheatedHost) Guess(guess string) (result string, finished bool) {
	guess = strings.ToUpper(guess)

	groups := make(map[string][]string)

	// Ensure guess is the same length as the answer
	if len(guess) != 5 { // TODO: Hardcoded for now
		return "Invalid Word length!", false
	}

	w.chances--

	// Caculate Hit and Presents
	for _, word := range w.WordList {
		resultRunes := make([]rune, len(guess))

		for i, r := range guess {
			letter := rune(word[i])
			if letter == r { // correct letter
				resultRunes[i] = 'O'

				continue
			}
			if _, present := w.wordLookupSet[word][r]; present {
				resultRunes[i] = '?'
			} else {
				resultRunes[i] = '_'
			}
		}

		pattern := string(resultRunes)
		groups[pattern] = append(groups[pattern], word)
	}

	var worstPattern string
	var largestGroup []string

	minHits := math.MaxInt32
	minPresents := 0

	for pattern, group := range groups {
		currHits := 0
		currPresents := 0

		for _, w := range pattern {
			if w == 'O' {
				currHits++
			} else if w == '?' {
				currPresents++
			}
		}
		if currHits < minHits {
			minHits = currHits
			minPresents = currPresents
			worstPattern = pattern
			largestGroup = group
		} else if currHits == minHits && currPresents < minPresents {
			minHits = currHits
			minPresents = currPresents
			worstPattern = pattern
			largestGroup = group
		}
	}
	w.WordList = largestGroup

	if worstPattern == "OOOOO" {
		return fmt.Sprintf("bingo! the correct answer is %s", w.WordList[0]), true
	}

	// Check if the game is over due to running out of chances
	if w.chances <= 0 {
		return "GAME OVER!!!", true
	}

	return worstPattern, finished
}
