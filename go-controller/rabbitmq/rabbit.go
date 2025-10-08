
package rabbitmq

import (
	"encoding/json"
	"log"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"os"
	"fmt"
)

// os.Getenv("HOST_QUEUE")
// var HOST_QUEUE = "10.24.191.38"
// var HOST_QUEUE string = os.Getenv("HOST_QUEUE")
// var HOST_QUEUE_PORT string = os.Getenv("HOST_QUEUE_PORT")

func Dial() (*amqp.Connection, error) {
	// return amqp.Dial("amqp://guest:guest@" + HOST_QUEUE + ":`HOST_QUEUE_PORT`/")
	host := os.Getenv("HOST_QUEUE")
    port := os.Getenv("HOST_QUEUE_PORT")
	// fmt.Println(HOST_QUEUE, HOST_QUEUE_PORT)
	// return amqp.Dial(fmt.Sprintf("amqp://user:password@%s:%s/", HOST_QUEUE, HOST_QUEUE_PORT))
	return amqp.Dial(fmt.Sprintf("amqp://user:password@%s:%s/", host, port))
}

func SendToQueue(runnerID string, data map[string]interface{}) string {
	msgID := uuid.New().String()
	return SendToQueueWithCustomID(runnerID, data, msgID)
}

// SendToQueueWithCustomID allows sending with a custom message ID
func SendToQueueWithCustomID(runnerID string, data map[string]interface{}, customMsgID string) string {
	host := os.Getenv("HOST_QUEUE")
	port := os.Getenv("HOST_QUEUE_PORT")
	conn, err := amqp.Dial(fmt.Sprintf("amqp://user:password@%s:%s/", host, port))
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
	
	// Sử dụng custom message ID thay vì generate mới
	data["id"] = customMsgID
	// data["reply_to"] = runnerID
	data["send_to"] = runnerID
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
		return ""
	}
	
	// log.Printf("Sent message to runner %s with ID: %s", runnerID, customMsgID)
	return customMsgID
}
