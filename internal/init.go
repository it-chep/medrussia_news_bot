package internal

import (
	"context"
	"log/slog"
	"medrussia_news_bot/internal/config"
	"medrussia_news_bot/internal/infrastructure/controller"
	"medrussia_news_bot/internal/infrastructure/controller/bot_controller"
	"medrussia_news_bot/internal/pkg/telegram"
	"os"
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
	a.controllers.botController = bot_controller.NewTelegramWebhookController(a.config, a.logger, a.bot)
	return a
}

func (a *App) initRestController(_ context.Context) *App {
	a.controllers.restController = controller.NewRestController(a.config, a.logger, a.controllers.botController)
	return a
}
