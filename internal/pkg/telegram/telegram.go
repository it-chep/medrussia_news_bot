package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"medrussia_news_bot/internal/config"
)

type Bot struct {
	Bot    *tgbotapi.BotAPI
	logger *slog.Logger
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
	return &Bot{
		Bot:    bot,
		logger: logger,
	}
}

func (bot *Bot) SendMessage(msg tgbotapi.MessageConfig) {
	_, err := bot.Bot.Send(msg)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("%s: Bot SendMessage", err))
	}

	return
}

func (bot *Bot) SendMessageAndGetId(msg tgbotapi.MessageConfig) int {
	sentMessage, err := bot.Bot.Send(msg)
	if err != nil {
		bot.logger.Error(fmt.Sprintf("%s: Bot SendMessageAndGetId", err))
	}
	return sentMessage.MessageID
}
