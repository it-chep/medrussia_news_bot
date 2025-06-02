package main

import (
	"context"
	"fmt"
	"medrussia_news_bot/internal"

	tele "gopkg.in/telebot.v3"
)

func main() {
	ctx := context.Background()
	fmt.Println("Старт сервера")
	app := internal.NewApp(ctx)
	app.Init()

	app.TgBot.Handle(tele.OnVoice, app.BotSvc.VoiceMess)

	app.TgBot.Handle(tele.OnText, app.BotSvc.OnText)

	app.TgBot.Start()
}
