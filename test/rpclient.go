package main

import (
	"github.com/streadway/amqp"
	"math/rand"
)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func Connect() (*amqp.Connection, *amqp.Channel) {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()

	return conn, ch
}

func Send() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	defer conn.Close()
	defer ch.Close()

	headers := map[string]interface{}{
		"1": "uuid1",
		"3": "uuid3",
		"4": "uuid3",
	}

	ch.ExchangeDeclare(
		"TestHeaders",
		"headers",
		false,
		false, //autoDelete
		false, //internal
		false, //noWait
		nil)


	ch.Publish(
		"TestHeaders",  // exchange
		"", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			Headers: headers,
			ContentType:   "text/plain",
			Body:          []byte("test"),
	})

}


func main() {
	Send()
}

