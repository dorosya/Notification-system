package main

import (
	"context"
	"fmt"
	"log"
	"notification-system/Notification_service/telegram"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

func notification_handler(ctx context.Context, b *bot.Bot, notifyContent string) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    841015314,
		Text:      fmt.Sprintf("Новое уведомление: %s", notifyContent),
		ParseMode: models.ParseModeMarkdown,
	})
}
func main() {

	//Сетап бота
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	tgbot := telegram.Bot_Init(ctx, os.Getenv("BOT_TOKEN"))

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panic(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		"Notifications",
		"",   // consumer tag
		true, // auto-ack, для моего проекта оставлю так, но для ретраев при стукании в тгшку лучшее конечно false  делать и добавлять доп обработку
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for msg := range msgs {
			go notification_handler(ctx, tgbot, string(msg.Body))
		}
	}()

	log.Println("Notification service started")

	<-ctx.Done()
}
