
package rabbitmq

import (
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"os"
	"fmt"
)


var HOST_QUEUE = os.Getenv("HOST_QUEUE")
var HOST_QUEUE_PORT = os.Getenv("HOST_QUEUE_PORT")
func Dial() (*amqp.Connection, error) {
	// return amqp.Dial("amqp://guest:guest@" + HOST_QUEUE + ":`HOST_QUEUE_PORT`/")
	return amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%s/", HOST_QUEUE, HOST_QUEUE_PORT))
}

func SendToQueue(runnerID string, data map[string]interface{}) string {
	// conn, err := amqp.Dial("amqp://guest:guest@" + HOST_QUEUE + ":5672/")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:%s/", HOST_QUEUE, HOST_QUEUE_PORT))
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
