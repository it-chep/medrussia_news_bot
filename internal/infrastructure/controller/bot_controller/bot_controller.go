package bot_controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log/slog"
	"medrussia_news_bot/internal/config"
	"medrussia_news_bot/internal/infrastructure/controller/bot_controller/dto"
	"medrussia_news_bot/internal/pkg/telegram"
	"net/http"
)

type TelegramWebhookController struct {
	cfg    *config.Config
	logger *slog.Logger
	bot    *telegram.Bot
}

func NewTelegramWebhookController(
	cfg *config.Config,
	logger *slog.Logger,
	bot *telegram.Bot,
) TelegramWebhookController {

	return TelegramWebhookController{
		cfg:    cfg,
		logger: logger,
		bot:    bot,
	}
}

func (t TelegramWebhookController) BotWebhookHandler(c *gin.Context) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			t.logger.Error(fmt.Sprintf("%s", err))
		}
	}(c.Request.Body)

	var update tgbotapi.Update
	if err := c.ShouldBindJSON(&update); err != nil {
		t.logger.Error("Error binding JSON", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	tgUser := t.getUserFromWebhook(update)
	tgMessage := t.getMessageFromWebhook(update)

	// Сначала проверяем на команду, потом на текстовое сообщение, потом callback
	if update.Message != nil {
		ctx := context.WithValue(context.Background(), "userID", update.Message.From.ID)
		if update.Message.IsCommand() {
			t.ForkCommands(ctx, update, tgUser, tgMessage)
		} else {
			t.ForkMessages(ctx, tgUser, tgMessage)
		}
	} else if update.CallbackQuery != nil {
		ctx := context.WithValue(context.Background(), "userID", update.CallbackQuery.From.ID)
		t.ForkCallbacks(ctx, update, tgUser, tgMessage)
	} else {
		t.logger.Warn(fmt.Sprintf("Unhandled update type: %+v", update))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
	return
}

func (t TelegramWebhookController) ForkCommands(ctx context.Context, update tgbotapi.Update, tgUser dto.TgUserDTO, tgMessage dto.MessageDTO) {

	switch update.Message.Command() {
	case "start":
		//...
	}
}

func (t TelegramWebhookController) ForkMessages(ctx context.Context, tgUser dto.TgUserDTO, tgMessage dto.MessageDTO) {

}

func (t TelegramWebhookController) ForkCallbacks(ctx context.Context, update tgbotapi.Update, tgUser dto.TgUserDTO, tgMessage dto.MessageDTO) {
}

func (t TelegramWebhookController) getUserFromWebhook(update tgbotapi.Update) dto.TgUserDTO {
	var tgUser dto.TgUserDTO
	var userJSON []byte
	var err error

	if update.CallbackQuery != nil {
		userJSON, err = json.Marshal(update.CallbackQuery.From)
	} else {
		userJSON, err = json.Marshal(update.Message.From)
	}

	if err != nil {
		t.logger.Error(fmt.Sprintf("Error marshaling user to JSON: %s", err))
		return dto.TgUserDTO{}
	}

	if err = json.Unmarshal(userJSON, &tgUser); err != nil {
		t.logger.Error(fmt.Sprintf("Error decoding JSON: %s", err))
		return dto.TgUserDTO{}
	}

	return tgUser
}

func (t TelegramWebhookController) getMessageFromWebhook(update tgbotapi.Update) dto.MessageDTO {
	var tgMessage dto.MessageDTO
	var userJSON []byte
	var err error

	if update.CallbackQuery != nil {
		userJSON, err = json.Marshal(update.CallbackQuery.Message)
	} else {
		userJSON, err = json.Marshal(update.Message)
	}

	if err != nil {
		t.logger.Error(fmt.Sprintf("Error marshaling user to JSON: %s", err))
		return dto.MessageDTO{}
	}

	if err = json.Unmarshal(userJSON, &tgMessage); err != nil {
		t.logger.Error(fmt.Sprintf("Error decoding JSON: %s", err))
		return dto.MessageDTO{}
	}

	return tgMessage
}
