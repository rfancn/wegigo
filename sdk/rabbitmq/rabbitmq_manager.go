package rabbitmq

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

type RabbitMQManager struct {
	Conn *amqp.Connection
	Ch *amqp.Channel
}

//NewRabbitMQManager: new rabbitmq manager
func NewRabbitMQManager(address string, port int) *RabbitMQManager {
	connStr := fmt.Sprintf("amqp://guest:guest@%s:%d/", address, port)
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Println("NewRabbitMQManager(): Error connect to RabbitMQ server:", err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println("NewRabbitMQManager(): Error open RabbitMQ channel:", err)
		return nil
	}

	return &RabbitMQManager{conn, ch}
}

//close rabbitmq channel and connections
func (m *RabbitMQManager) Close() {
	if m.Ch != nil {
		err := m.Ch.Close()
		if err != nil {
			log.Println("RabbitMQManager Close(): Error close RabbitMQ channel:", err)
		}
	}

	if m.Conn != nil {
		m.Conn.Close()
	}
}
