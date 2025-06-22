package internal

import (
	"context"
	"log"
	"log/slog"
	"medrussia_news_bot/internal/config"
	"medrussia_news_bot/internal/infrastructure/controller"
	"medrussia_news_bot/internal/infrastructure/controller/bot_controller"
	"medrussia_news_bot/internal/infrastructure/repo"
	"medrussia_news_bot/internal/pkg/postgres"
	"medrussia_news_bot/internal/pkg/telegram"
	"net/http"
	"os"
	"time"
)

func (a *App) initConfig(_ context.Context) *App {
	a.config = config.NewConfig()
	return a
}

func (a *App) initLogger(_ context.Context) *App {
	a.logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	return a
}

func (a *App) initBot(_ context.Context) *App {
	a.bot = telegram.NewTelegramBot(a.config, a.logger)
	return a
}

func (a *App) initBotController(_ context.Context) *App {
	a.controllers.botController = bot_controller.NewTelegramWebhookController(a.config, a.logger, a.bot, a.repo)
	return a
}

func (a *App) initPgxConn(ctx context.Context) *App {
	client, err := postgres.NewClient(ctx, a.config.StorageConfig)
	if err != nil {
		log.Fatal(err)
	}
	a.pgxClient = client
	a.logger.Info("init pgxclient", a.pgxClient)
	return a
}

func (a *App) initRepo(_ context.Context) *App {
	a.repo = repo.NewRepo(a.pgxClient, a.logger)
	return a
}

func (a *App) iniControllers(_ context.Context) *App {
	a.controllers = controllers{}
	return a
}

func (a *App) initServer(_ context.Context) *App {
	a.controllers.restController = controller.NewRestController(a.config, a.logger, a.controllers.botController)

	a.server = &http.Server{
		Addr:         a.config.HTTPServer.Address,
		Handler:      a.controllers.restController,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	return a
}
