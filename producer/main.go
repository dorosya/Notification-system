package main

import (
	"log"
	"net/http"
	"notification-system/producer/net/http/handlers"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func rabbitCon() *amqp.Connection {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	return conn
}

func routerSetup(conn *amqp.Connection) {
	router := gin.Default()
	router.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})
	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	handler := handlers.EventHandler{Conn: conn}
	router.POST("/event", handler.EventsHandler)
	router.Run(":8080")
}

func main() {
	conn := rabbitCon()
	defer conn.Close()

	routerSetup(conn)
}
