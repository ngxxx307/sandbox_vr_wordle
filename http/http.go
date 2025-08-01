package main

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/ngxxx307/sandbox_vr_wordle/config"
	"github.com/ngxxx307/sandbox_vr_wordle/controller"
	"github.com/ngxxx307/sandbox_vr_wordle/hub"
	"github.com/ngxxx307/sandbox_vr_wordle/routes"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(config.NewConfig),

		fx.Provide(controller.NewWebsocket),
		fx.Provide(controller.NewGameLoungeController),
		fx.Provide(controller.NewWordleController),

		fx.Provide(hub.NewHub),

		fx.Provide(NewEchoServer),
		fx.Invoke(routes.SetupWebSocketRoute),

		fx.Invoke(StartEchoServer),
		fx.Invoke(RunHub),
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

func RunHub(lc fx.Lifecycle, h *hub.Hub) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go h.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			h.Stop()
			return nil
		},
	})
}
