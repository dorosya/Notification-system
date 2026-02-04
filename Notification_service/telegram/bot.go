package telegram

import (
	"context"
	"log"

	"github.com/go-telegram/bot"
)

func Bot_Init(ctx context.Context, token string) *bot.Bot {
	b, err := bot.New(token)
	if nil != err {
		panic(err)
	}

	go b.Start(ctx)
	log.Println("Telegram bot started")
	return b
}
