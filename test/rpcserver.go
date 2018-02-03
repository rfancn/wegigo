package main

import (
	"log"
	"github.com/streadway/amqp"
)

/**
func getEnabledApps() map[string]string{

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   "http://127.0.0.4379",
		DialTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	resp, err := cli.Get(ctx, "/app/enabled")
	cancel()
	if err != nil {
		log.Println("EtcdManager Put(): Error read from etcd:", err)
		return nil
	}

	apps := make(map[string]string)
	if err := json.Unmarshal(resp.Kv.Value, &apps); err != nil {
		log.Printf("AppManager GetEnabledApps(): Error unmarshal map[string]string:%v", err)
		return nil
	}

	return apps
}
**/

func work1() {
	conn, _ := amqp.Dial("amqp://guest:guest@localhost:5672/")
	ch, _ := conn.Channel()
	defer conn.Close()
	defer ch.Close()

	headers := map[string]interface{}{
		"x-match": "all",
		"1": "uuid1",
		"3": "uuid3"}

	ch.ExchangeDeclare(
		"TestHeaders",
		"headers",
		false,
		false, //autoDelete
		false, //internal
		false, //noWait
		nil)

		q, _ := ch.QueueDeclare(
			"", // name
			false,       // durable
			false,       // delete when usused
			false,       // exclusive
			false,       // no-wait
			nil,         // arguments
		)

		ch.QueueBind(
			q.Name,     // queue name
			"",       // routing key
			"TestHeaders", // exchange name
			false,        //no-wait
			headers)           //arguments

	messages, _ := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	for m := range messages {
		log.Println("work1:", string(m.Body))
	}
}


func main() {
	work1()
}