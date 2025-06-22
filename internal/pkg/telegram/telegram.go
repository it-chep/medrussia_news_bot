package telegram

import (
	"fmt"
	"log/slog"
	"medrussia_news_bot/internal/config"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	Bot         *tgbotapi.BotAPI
	logger      *slog.Logger
	adminChatID int64
}

func NewTelegramBot(cfg *config.Config, logger *slog.Logger) *Bot {
	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	bot.Debug = true
	if err != nil {
		panic("can't create bot instance")
	}

	wh, _ := tgbotapi.NewWebhook(cfg.Bot.WebhookURL + bot.Token + "/")
	_, err = bot.Request(wh)
	if err != nil {
		panic("can't while request set webhook")
	}

	_, err = bot.GetWebhookInfo()

	if err != nil {
		panic("error while getting webhook")
	}

	adminChatID, err := strconv.ParseInt(cfg.Bot.AdminChatID, 10, 64)
	if err != nil {
		panic("Неправильно указан adminChatID")
	}
	logger.Info("Telegram bot initialized")
	return &Bot{
		Bot:         bot,
		logger:      logger,
		adminChatID: adminChatID,
	}
}

func (bot *Bot) SendMessage(msg tgbotapi.MessageConfig) {
	_, err := bot.Bot.Send(msg)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("%s: Bot SendMessage", err))
	}

	return
}

func (bot *Bot) ForwardMessage(msg *tgbotapi.Message) {
	// Создаем конфиг для пересылки сообщения
	forwardConfig := tgbotapi.NewForward(
		bot.adminChatID,
		msg.Chat.ID,
		msg.MessageID,
	)

	// Отправляем сообщение
	_, err := bot.Bot.Send(forwardConfig)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("failed to forward message: %v", err))
	}

	return
}

func (bot *Bot) ForwardAdminMessageToUser(chatID int64, msg *tgbotapi.Message) int64 {
	forwardConfig := tgbotapi.NewMessage(
		chatID,
		fmt.Sprintf("Ответ администрации бота: \n\n %s", msg.Text),
	)

	// Отправляем сообщение
	_, err := bot.Bot.Send(forwardConfig)
	if err != nil {
		_, _ = bot.Bot.Send(tgbotapi.NewMessage(bot.adminChatID, "Ошибка при отправке сообщения"))
		bot.logger.Error(fmt.Sprintf("failed to forward message: %v", err))
	}

	return int64(msg.MessageID)
}

func (bot *Bot) SendMessageToAdmin(
	text string,
) (messageID int64, err error) {
	msg := tgbotapi.NewMessage(bot.adminChatID, text)

	message, err := bot.Bot.Send(msg)
	if err != nil {
		bot.logger.Error("Ошибка пересылки сообщения: " + err.Error())
	}

	return int64(message.MessageID), err
}

func (bot *Bot) ForwardMessageToAdminChat(
	replyToMessageID,
	chatID int64,
	text string,
) (messageID int64, err error) {
	msg := tgbotapi.NewMessage(bot.adminChatID, text)

	if replyToMessageID != 0 {
		msg.ReplyToMessageID = int(replyToMessageID)
	}

	// callback читается как "'действие' _ 'ID пользователя'
	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
			"Закрыть обращение ❇️",
			fmt.Sprintf("close_%d", chatID),
		)),
	)

	message, err := bot.Bot.Send(msg)
	if err != nil {
		bot.logger.Error("Ошибка пересылки сообщения: " + err.Error())
	}

	bot.logger.Info(fmt.Sprintf("Сообщение пользователя: %v", message))

	return int64(message.MessageID), err
}

// CleanMessageButtonsInAdminChat убирает inline кнопки в сообщении по его ID
func (bot *Bot) CleanMessageButtonsInAdminChat(
	messageID int64,
) (editedMessageID int64, err error) {
	emptyMarkup := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{}}

	editMarkup := tgbotapi.NewEditMessageReplyMarkup(bot.adminChatID, int(messageID), emptyMarkup)

	message, err := bot.Bot.Send(editMarkup)
	if err != nil {
		bot.logger.Error("Ошибка редактирования сообщения в чате админов: " + err.Error())
	}

	return int64(message.MessageID), err
}

func (bot *Bot) SetCloseButtonInAdminChat(
	messageID int64,
) (editedMessageID int64, err error) {
	if messageID == 0 {
		return
	}

	closedMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(
			"⭕️ Обращение закрыто ⭕️️",
			"ignore",
		)),
	)

	editMarkup := tgbotapi.NewEditMessageReplyMarkup(bot.adminChatID, int(messageID), closedMarkup)

	message, err := bot.Bot.Send(editMarkup)
	if err != nil {
		bot.logger.Error("Ошибка при изменение кнопки на закрытие в чате админов: " + err.Error())
	}

	return int64(message.MessageID), err
}

func (bot *Bot) SendMessageAndGetId(msg tgbotapi.MessageConfig) int {
	sentMessage, err := bot.Bot.Send(msg)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("%s: Bot SendMessageAndGetId", err))
	}
	return sentMessage.MessageID
}
