package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	"github.com/ngxxx307/sandbox_vr_wordle/matchMaker"
)

// GameContext provides type-safe access to dependencies
type GameContext struct {
	EchoCtx    echo.Context
	Config     *config.Config
	MatchMaker *matchMaker.MatchMaker
}

// NewGameContext creates a new game context with dependencies
func NewGameContext(echoCtx echo.Context, cfg *config.Config, matcher *matchMaker.MatchMaker) *GameContext {
	return &GameContext{
		EchoCtx:    echoCtx,
		Config:     cfg,
		MatchMaker: matcher,
	}
}
