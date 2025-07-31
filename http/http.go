package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	"github.com/ngxxx307/sandbox_vr_wordle/routes"
	"github.com/ngxxx307/sandbox_vr_wordle/websocket"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(config.NewConfig),

		fx.Provide(websocket.NewWebsocket),

		fx.Provide(NewEchoServer),
		fx.Invoke(routes.SetupWebSocketRoute),
		fx.Invoke(StartEchoServer),
	).Run()
}

func NewEchoServer(lc fx.Lifecycle) *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	e.Pre(middleware.AddTrailingSlash())
	e.Use(middleware.CORS())
	return e
}

func StartEchoServer(lc fx.Lifecycle, e *echo.Echo, cfg *config.Config) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			addr := fmt.Sprintf("0.0.0.0:%s", cfg.ServerPort)
			fmt.Printf("Server listening on port %s\n", cfg.ServerPort)
			go e.Start(addr)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Graceful Shutdown")
			return e.Shutdown(ctx)
		},
	})
}
