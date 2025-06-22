package internal

import (
	"context"
	"log/slog"
	"medrussia_news_bot/internal/config"
	"medrussia_news_bot/internal/infrastructure/controller"
	"medrussia_news_bot/internal/infrastructure/controller/bot_controller"
	"medrussia_news_bot/internal/infrastructure/repo"
	"medrussia_news_bot/internal/pkg/postgres"
	"medrussia_news_bot/internal/pkg/telegram"
	"net/http"
)

type controllers struct {
	botController  bot_controller.TelegramWebhookController
	restController *controller.RestController
}

type App struct {
	logger      *slog.Logger
	config      *config.Config
	pgxClient   postgres.Client
	controller  controllers
	bot         *telegram.Bot
	controllers controllers
	server      *http.Server
	repo        *repo.Repo
}

func NewApp(ctx context.Context) *App {
	a := &App{}

	a.initLogger(ctx).
		initConfig(ctx).
		initPgxConn(ctx).
		initRepo(ctx).
		initBot(ctx).
		iniControllers(ctx).
		initBotController(ctx).
		initServer(ctx)

	return a
}

func (a *App) Run(_ context.Context) error {
	a.logger.Info("start server")
	return a.server.ListenAndServe()
}
