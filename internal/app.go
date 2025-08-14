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

	//if !a.config.UseWebhook() {
	//	// Режим поллинга
	//	for update := range a.bot.GetUpdates() {
	//		go func() {
	//			if update.FromChat() == nil || update.SentFrom() == nil {
	//				return
	//			}
	//
	//			mediaID := ""
	//			mediaType := commonDto.Unknown
	//
	//			if update.Message.Video != nil {
	//				mediaID = update.Message.Video.FileID
	//				mediaType = commonDto.Video
	//			}
	//
	//			if len(update.Message.Photo) != 0 {
	//				mediaID = update.Message.Photo[0].FileID
	//				mediaType = commonDto.Photo
	//			}
	//
	//			if update.Message.Document != nil {
	//				mediaID = update.Message.Document.FileID
	//				mediaType = commonDto.Document
	//			}
	//
	//			if update.Message.Audio != nil {
	//				mediaID = update.Message.Audio.FileID
	//				mediaType = commonDto.Audio
	//			}
	//
	//			if update.Message.Voice != nil {
	//				mediaID = update.Message.Voice.FileID
	//				mediaType = commonDto.Voice
	//			}
	//
	//			if update.Message.VideoNote != nil {
	//				mediaID = update.Message.VideoNote.FileID
	//				mediaType = commonDto.VideoNote
	//			}
	//			fmt.Print("Обработка ивента", update.Message.Chat.ID, update)
	//			a.controllers.botController.ForkMessages(
	//				context.Background(), update, dto.TgUserDTO{}, dto.MessageDTO{
	//					MediaID:   mediaID,
	//					MediaType: mediaType,
	//					Chat:      &dto.Chat{ID: update.Message.Chat.ID},
	//				},
	//			)
	//
	//		}()
	//	}
	//}
	return a.server.ListenAndServe()
}
