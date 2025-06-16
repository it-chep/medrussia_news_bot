package internal

import (
	"context"
	"log/slog"
	"medrussia_news_bot/internal/config"
	"medrussia_news_bot/internal/infrastructure/controller"
	"medrussia_news_bot/internal/infrastructure/controller/bot_controller"
	"medrussia_news_bot/internal/pkg/telegram"
	"net/http"
)

type controllers struct {
	botController  bot_controller.TelegramWebhookController
	restController controller.RestController
}

type App struct {
	logger      *slog.Logger
	config      *config.Config
	controller  controllers
	bot         *telegram.Bot
	controllers controllers
	server      *http.Server
}

func NewApp(ctx context.Context) *App {
	a := &App{}

	a.initLogger(ctx).
		initConfig(ctx).
		initBot(ctx).
		initBotController(ctx).
		initRestController(ctx)

	return a
}

func (app *App) Run(_ context.Context) error {
	app.logger.Info("start server")
	return app.server.ListenAndServe()
}
