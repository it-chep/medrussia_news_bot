package bot_controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"medrussia_news_bot/internal/config"
	"medrussia_news_bot/internal/infrastructure/controller/bot_controller/dto"
	"medrussia_news_bot/internal/infrastructure/repo"
	"medrussia_news_bot/internal/pkg/telegram"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Repo interface {
	GetUser(ctx context.Context, userID int64) (*repo.UserDialog, error)
	CreateUser(ctx context.Context, userID int64) error
	UpdateLastAdminMessage(ctx context.Context, userID, messageID int64) error
	UpdateLastUserMessage(ctx context.Context, userID, messageID int64) error
	TurnOnAvailable(ctx context.Context, userID int64) error
	CloseAppeal(ctx context.Context, userID int64) error
}

const (
	startMessage = `
Здравствуйте!

Вы связались с редакцией «Медицинская Россия».

Мы принимаем материалы, фото, видео и истории, связанные с проблемами медработников в сфере здравоохранения. Вы можете рассказать нам о:

— трудовых конфликтах и нарушениях прав ваших коллег,
— проблемах в медицинском образовании и подготовке кадров,
— несправедливом преследовании врачей,
— коррупции и неэффективном управлении в системе здравоохранения,
— нехватке оборудования, кадров или лекарств,
— угрозах врачам,
— и других острых или системных вопросах, требующих внимания.
 
Присылайте свои обращения прямо сюда в этот диалог и в ближайшее время вам ответят.

Обратите внимание: мы не гарантируем публикацию статьи по каждому обращению, но обязательно постараемся внимательно изучить вашу проблему.
`
)

// TelegramWebhookController контроллер бота
type TelegramWebhookController struct {
	cfg    *config.Config
	logger *slog.Logger
	bot    *telegram.Bot
	repo   Repo
}

// NewTelegramWebhookController конструктор
func NewTelegramWebhookController(
	cfg *config.Config,
	logger *slog.Logger,
	bot *telegram.Bot,
	repo Repo,
) TelegramWebhookController {
	return TelegramWebhookController{
		cfg:    cfg,
		logger: logger,
		bot:    bot,
		repo:   repo,
	}
}

// BotWebhookHandler хендлер реагирующий на все вебхуки бота
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
			// Если чат админский
			if strconv.FormatInt(update.Message.Chat.ID, 10) == t.cfg.Bot.AdminChatID {
				t.ForkAdminMessage(ctx, update)
				return
			}
			t.ForkMessages(ctx, update, tgUser, tgMessage)
		}
	} else if update.CallbackQuery != nil {
		ctx := context.WithValue(context.Background(), "userID", update.CallbackQuery.From.ID)
		t.ForkCallbacks(ctx, update)
	} else {
		t.logger.Warn(fmt.Sprintf("Unhandled update type: %+v", update))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
	return
}

// ForkCommands обработка всех сообщений типа Command
func (t TelegramWebhookController) ForkCommands(ctx context.Context, update tgbotapi.Update, tgUser dto.TgUserDTO, tgMessage dto.MessageDTO) {

	switch update.Message.Command() {
	case "start":
		err := t.repo.CreateUser(ctx, update.Message.From.ID)
		if err != nil {
			t.logger.Error(fmt.Sprintf("%s", err))
			return
		}
		t.bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, startMessage))
	}
}

// ForkMessages обработка всех сообщений типа MESSAGE
func (t TelegramWebhookController) ForkMessages(ctx context.Context, update tgbotapi.Update, tgUser dto.TgUserDTO, tgMessage dto.MessageDTO) {
	user, err := t.repo.GetUser(ctx, update.Message.From.ID)
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s", err))

		_ = t.repo.CreateUser(ctx, update.Message.From.ID)
		user, err = t.repo.GetUser(ctx, update.Message.From.ID)
	}

	if !user.Available {
		// отвечаем пользователю
		t.bot.SendMessage(tgbotapi.NewMessage(update.Message.Chat.ID, "Пожалуйста ожидайте, вам скоро ответят"))
	}
	// ставим флаг активности диалога
	t.setActiveUserFlag(ctx, user)
	// шлем админам
	t.forwardToAdmin(ctx, user, update)
}

// ForkAdminMessage пересылка пользователю сообщение админа
func (t TelegramWebhookController) ForkAdminMessage(ctx context.Context, update tgbotapi.Update) {
	replyTo := update.Message.ReplyToMessage
	if replyTo == nil {
		return
	}
	markup := replyTo.ReplyMarkup
	if markup == nil {
		return
	}

	userID := t.parseReplyCloseKeyboard(update.Message.ReplyToMessage)
	if userID == 0 {
		return
	}

	user, err := t.repo.GetUser(ctx, userID)
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s", err))
		return
	}

	forwardMessageID := t.bot.ForwardAdminMessageToUser(user.UserID, update.Message)

	t.saveAdminMessageIDToUser(ctx, user, forwardMessageID)
}

// saveAdminMessageIDToUser сохранение последнего сообщения админа пользователю
func (t TelegramWebhookController) saveAdminMessageIDToUser(ctx context.Context, user *repo.UserDialog, adminMessageID int64) {
	err := t.repo.UpdateLastAdminMessage(ctx, user.UserID, adminMessageID)
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s", err))
		return
	}
}

// saveUserMessageID сохранение последнего сообщения пользователя
func (t TelegramWebhookController) saveUserMessageID(ctx context.Context, user *repo.UserDialog, update tgbotapi.Update) {
	err := t.repo.UpdateLastUserMessage(ctx, user.UserID, int64(update.Message.MessageID))
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s", err))
		return
	}
}

// parseReplyCloseKeyboard парсинг клавиатуры при ответе админа пользователю
func (t TelegramWebhookController) parseReplyCloseKeyboard(replyToMessage *tgbotapi.Message) (userID int64) {
	for _, line := range replyToMessage.ReplyMarkup.InlineKeyboard {
		for _, button := range line {
			var callbackData string

			if button.CallbackData != nil {
				callbackData = *button.CallbackData
			}

			if strings.Contains(callbackData, "close") {
				return t.parseCloseCallback(callbackData)
			}
		}
	}

	return userID
}

// ForkCallbacks Обработка колбека сообщения
func (t TelegramWebhookController) ForkCallbacks(ctx context.Context, update tgbotapi.Update) {
	t.processCloseCallback(ctx, update)
}

// getUserFromWebhook получение пользователя из вебхука
func (t TelegramWebhookController) getUserFromWebhook(update tgbotapi.Update) dto.TgUserDTO {
	var tgUser dto.TgUserDTO
	var userJSON []byte
	var err error

	if update.CallbackQuery != nil {
		userJSON, err = json.Marshal(update.CallbackQuery.From)
	} else if update.Message != nil {
		userJSON, err = json.Marshal(update.Message.From)
	} else {
		t.logger.Error("Cannot get user from webhook - no valid user data found", update)
		return dto.TgUserDTO{}
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

// setActiveUserFlag ставим флаг активности диалога
func (t TelegramWebhookController) setActiveUserFlag(ctx context.Context, user *repo.UserDialog) {
	if user.Available {
		return
	}
	user.Available = true

	err := t.repo.TurnOnAvailable(ctx, user.UserID)
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s", err))
	}
}

// forwardToAdmin пересылаем админам и ставим пользователю id сообщения в админ чате
func (t TelegramWebhookController) forwardToAdmin(ctx context.Context, user *repo.UserDialog, update tgbotapi.Update) {
	text := fmt.Sprintf(
		"Пользователь: @%s\nИмя: %s %s\n\nТекст сообщения: %s",
		update.Message.Chat.UserName, update.Message.Chat.LastName, update.Message.Chat.FirstName, update.Message.Text,
	)
	forwardMessageID, err := t.bot.ForwardMessageToAdminChat(user.LastAdminMessageID.Int64, update.Message.Chat.ID, text)
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s", err))
		return
	}
	err = t.repo.UpdateLastUserMessage(ctx, user.UserID, forwardMessageID)
	if err != nil {
		t.logger.Error(fmt.Sprintf("%s", err))
		return
	}

	_, err = t.bot.CleanMessageButtonsInAdminChat(user.LastUserMessageID.Int64)
	if err != nil {
		t.logger.Error("Ошибка при очистки клавиатуры сообщения администратора: " + err.Error())
	}
}

// getMessageFromWebhook получение сообщения из вебхука
func (t TelegramWebhookController) getMessageFromWebhook(update tgbotapi.Update) dto.MessageDTO {
	var tgMessage dto.MessageDTO
	var userJSON []byte
	var err error

	if update.CallbackQuery != nil {
		userJSON, err = json.Marshal(update.CallbackQuery.Message)
	} else if update.Message != nil {
		userJSON, err = json.Marshal(update.Message)
	} else {
		t.logger.Error("Cannot get user from webhook - no valid user data found", update)
		return dto.MessageDTO{}
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

// parseCloseCallback парсинг колбека о закрытии обращения в админке
func (t TelegramWebhookController) parseCloseCallback(callbackData string) int64 {
	data := strings.Split(callbackData, "_")
	if len(data) < 2 {
		return 0
	}

	userIDFromCallbackData, err := strconv.Atoi(data[1])
	userID := int64(userIDFromCallbackData)

	if err != nil {
		t.logger.Error("Ошибка при парсинге callback, не найден UserID" + err.Error())
	}

	return userID
}

// processCloseCallback обработка сallback закрытия обращения пользователя со стороны админа
func (t TelegramWebhookController) processCloseCallback(ctx context.Context, update tgbotapi.Update) {
	userID := t.parseCloseCallback(update.CallbackData())
	if userID == 0 {
		return
	}

	messageID := update.CallbackQuery.Message.MessageID

	updatedMessageID, err := t.bot.SetCloseButtonInAdminChat(int64(messageID))
	if err != nil {
		t.notifyAdmin(
			fmt.Sprintf(
				"Ошибка при закрытии обращения (изменение сообщения messageID: %d) %v",
				updatedMessageID,
				err.Error(),
			),
		)

		return
	}

	err = t.repo.CloseAppeal(ctx, userID)
	if err != nil {
		t.logger.Error("Ошибка при закрытии обращения пользователя: " + err.Error())
	}
}

// notifyAdmin оповещение админа о чем-то
func (t TelegramWebhookController) notifyAdmin(message string) {
	_, err := t.bot.SendMessageToAdmin(message)
	if err != nil {
		t.logger.Error(err.Error())
	}
}
