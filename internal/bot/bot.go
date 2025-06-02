package bot

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/samber/lo"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	tgBot *tele.Bot
}

func New(tgBot *tele.Bot) *Bot {
	return &Bot{tgBot: tgBot}
}

func (b *Bot) VoiceMess(c tele.Context) error {
	fmt.Println("дернулась ручка VoiceMess")
	ctx := context.Background()
	commands, err := b.repo.GetCommands(ctx, c.Chat().ID)
	if err != nil {
		return err
	}
	if len(commands) == 0 {
		return c.Send("У вас нет ни одного сценария")
	}

	voice := c.Message().Voice

	file, err := b.tgBot.FileByID(voice.FileID)
	if err != nil {
		return c.Send("Не удалось получить файл голосового сообщения")
	}

	f, err := b.tgBot.File(&file)
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	matched, err := b.neuro.GetAudio(ctx, lo.Map(commands, func(command repo.Command, _ int) dto.Command {
		return dto.Command{Name: command.Command}
	}), bytes)
	if err != nil {
		return err
	}

	err = b.yandex.Match(ctx, lo.FindOrElse(commands, repo.Command{}, func(command repo.Command) bool {
		return command.Command == matched.Name
	}))
	if err != nil {
		return err
	}

	return c.Send(fmt.Sprintf("Выполнена команда \"%s\"", matched.Name))
}

func (b *Bot) Help(c tele.Context) error {
	fmt.Println("дернулась ручка Help")

	helpText := `Доступные команды:
/start - Начать работу с ботом
/adddevice - Добавить устройство
/createdommand - Добавить сценарий умного устройства
/exit - Очистка состояний
/mycommands - Мои сценарии
/deletecommand - Удалить сценарий`
	return c.Send(helpText)
}
