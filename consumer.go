package main

import (
	"github.com/streadway/amqp"
	"log"
	"time"
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

	// 1つのメッセージを処理し終わってから、次のメッセージを受け取るようにするための設定
	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer

		// ACKを送信せずにコンシューマーが死んだ（チャネルが閉じた、接続が閉じた、またはTCP接続が失われた）場合、RabbitMQはメッセージが完全に処理されなかったことを認識し、それを再キューイングします。
		// 同時に他の消費者がオンラインにいる場合、すぐに別の消費者に再配信します。
		// そうすれば、workerがときどき死んでも、メッセージが失われないことが確実になります。
		// つまり、auto-ackをtrueにすると、workerが死んだら、メッセージが失われるということです。
		// 参考リンク
		// * https://www.rabbitmq.com/tutorials/tutorial-two-go.html
		// * https://qiita.com/k0001/items/7ae49db0621afda88a80
		false, // auto-ack

		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			time.Sleep(3 * time.Second)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
