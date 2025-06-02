package internal

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"medrussia_news_bot/internal/bot"
	"time"

	"github.com/georgysavva/scany/v2/dbscan"
	"github.com/georgysavva/scany/v2/pgxscan"

	tele "gopkg.in/telebot.v3"
)

func init() {
	// ignore db columns that doesn't exist at the destination
	dbscanAPI, err := pgxscan.NewDBScanAPI(dbscan.WithAllowUnknownColumns(true))
	if err != nil {
		panic(err)
	}

	api, err := pgxscan.NewAPI(dbscanAPI)
	if err != nil {
		panic(err)
	}

	pgxscan.DefaultAPI = api
}

type App struct {
	BotSvc *bot.Bot
	TgBot  *tele.Bot
}

func NewApp(ctx context.Context) *App {
	fmt.Println("Инициализация бота")
	pref := tele.Settings{
		Token:  "",
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	botTg, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Инициализация постгреса")
	// pgx
	conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@postgres:5432/mirea?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	return &App{
		TgBot: botTg,
	}
}

func (a *App) Init() {
	fmt.Println("ИНИТ")

	a.BotSvc = bot.New(a.TgBot)
}
