package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EventHandler struct {
	Conn *amqp.Connection
}

func (h *EventHandler) EventsHandler(c *gin.Context) {
	Channel, err := h.Conn.Channel()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	defer Channel.Close()
	if err != nil {
		log.Panicf("Failed to declare a queue")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = Channel.QueueDeclare(
		"Notifications", // name
		true,            // durable (сохраняется после перезагрузки сервера)
		false,           // delete when unused
		false,           // exclusive (только для этого соединения)
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Panicf(err.Error())
	}
	err = Channel.PublishWithContext(ctx,
		"",              // exchange
		"Notifications", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(c.PostForm("message")),
		})
	if err != nil {
		log.Println("Failed to publish a message")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": c.PostForm("message")})
}
