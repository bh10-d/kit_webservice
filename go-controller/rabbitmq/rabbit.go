
package rabbitmq

import (
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)


var RABBITMQ_HOST = "192.168.174.128"
func Dial() (*amqp.Connection, error) {
	return amqp.Dial("amqp://guest:guest@" + RABBITMQ_HOST + ":5672/")
}

func SendToQueue(runnerID string, data map[string]interface{}) string {
	conn, err := amqp.Dial("amqp://guest:guest@" + RABBITMQ_HOST + ":5672/")
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return ""
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return ""
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(
		runnerID,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare queue: %v", err)
		return ""
	}
	msgID := uuid.New().String()
	data["id"] = msgID
	data["reply_to"] = runnerID
	body, _ := json.Marshal(data)
	err = ch.Publish(
		"",
		runnerID,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	}
	return msgID
}
