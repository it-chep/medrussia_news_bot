package main

import (
	"context"
	"log"
	"medrussia_news_bot/internal"
)

func main() {
	ctx := context.Background()
	log.Fatal(internal.NewApp(ctx).Run(ctx))
}
