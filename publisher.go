package main

import (
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// キューの宣言は冪等です
	// まだ存在しない場合にのみ作成されます
	// メッセージの内容はバイト配列なので、好きなものをエンコードできます
	q, err := ch.QueueDeclare(
		"SampleQueue", // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return
	}

	body := "Hello World!"
	for i := 0; i < 10; i++ {
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body + strconv.Itoa(i)),
			})
		if err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
			return
		}
	}
}
