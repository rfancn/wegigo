package rabbitmq

import (
	"log"
	"github.com/streadway/amqp"
)

type RabbitMQManager struct {
	Conn        *amqp.Connection
	mainChannel *amqp.Channel
}

//NewRabbitMQManager: new rabbitmq manager
func NewRabbitMQManager(url string) (*RabbitMQManager, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &RabbitMQManager{conn, ch}, nil
}

//
// durable: If set when creating a new exchange, the exchange will be marked as durable.
// Durable exchanges remain active when a server restarts.
// Non-durable exchanges (transient exchanges) are purged if/when a server restarts.
//The server MUST support both durable and transient exchanges.
//
//autoDelete: If set, the exchange is deleted when all queues have finished using it.
//
//internal: If set, the exchange may not be used directly by publishers, but only when bound to other exchanges.
//Internal exchanges are used to construct wiring that is not visible to applications.
//
//nowait: If set, the server will not respond to the method.
//The client should not wait for a reply method. If the server could not complete the method it will raise a channel or connection exception.
func (m *RabbitMQManager) DeclareTopicExchange(name string, durable bool) bool {
	err := m.mainChannel.ExchangeDeclare(
		name,
		"topic",
		durable,
		false, //autoDelete
		false, //internal
		false, //noWait
		nil)

	if err != nil {
		log.Println("RabbitMQManager DeclareTopicExchane(): Error declar topic exchane:", err)
		return false
	}
	return true
}

func (m *RabbitMQManager) DeclareDirectExchange(name string, durable bool) bool {
	err := m.mainChannel.ExchangeDeclare(
		name,
		"direct",
		durable,
		false, //autoDelete
		false, //internal
		false, //noWait
		nil)

	if err != nil {
		log.Println("RabbitMQManager DeclareDirectExchange(): Error declar direct exchane:", err)
		return false
	}
	return true
}

func (m *RabbitMQManager) DeclareHeadersExchange(name string) bool {
	err := m.mainChannel.ExchangeDeclare(
		name,
		"headers",
		false, //durable
		false, //autoDelete
		false, //internal
		false, //noWait
		nil)

	if err != nil {
		log.Println("RabbitMQManager DeclareDirectExchange(): Error declar direct exchane:", err)
		return false
	}
	return true
}

//durable: If set when creating a new queue, the queue will be marked as durable.
// Durable queues remain active when a server restarts. Non-durable queues (transient queues) are purged if/when a server restarts.
// Note that durable queues do not necessarily hold persistent messages, although it does not make sense to send persistent messages to a transient queue.
//The server MUST recreate the durable queue after a restart.
//The server MUST support both durable and transient queues.
//
//exclusive: Exclusive queues may only be accessed by the current connection, and are deleted when that connection closes.
//Passive declaration of an exclusive queue by other connections are not allowed.
//
//autoDelete: If set, the queue is deleted when all consumers have finished using it.
// The last consumer can be cancelled either explicitly or because its channel is closed.
// If there was no consumer ever on the queue, it won't be deleted.
// Applications can explicitly delete auto-delete queues using the Delete method as normal.
func (m *RabbitMQManager) DeclareQueue(name string, durable bool) string {
	q, err := m.mainChannel.QueueDeclare(
		name,    // name
		durable, // durable
		false, // autoDelete
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Println("RabbitMQManager DeclareQueue(): Error declare queue:", err)
		return ""
	}

	return q.Name
}

func (m *RabbitMQManager) DeclareTempQueue() string {
	q, err := m.mainChannel.QueueDeclare(
		"",    // name
		false, // durable
		true,  // autoDelete
		false,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		log.Println("RabbitMQManager DeclareQueue(): Error declare queue:", err)
		return ""
	}

	return q.Name
}

//This method binds a queue to an exchange.
// Until a queue is bound it will not receive any messages.
// In a classic messaging model, store-and-forward queues are bound to a direct exchange
// and subscription queues are bound to a topic exchange.
func (m *RabbitMQManager) BindQueue(queueName string, exchangeName string, routingKey string) bool {
	err := m.mainChannel.QueueBind(
		queueName,     // queue name
		routingKey,    // routing key
		exchangeName,  // exchange name
		false,        //no-wait
		nil)           //arguments

	if err != nil {
		log.Println("RabbitMQManager BindQueue(): Error bind queue:to exchange", err)
		return false
	}
	return true
}

//This method binds a queue to an exchange.
// Until a queue is bound it will not receive any messages.
// In a classic messaging model, store-and-forward queues are bound to a direct exchange
// and subscription queues are bound to a topic exchange.
func (m *RabbitMQManager) BindQueueWithHeaders(queueName string, exchangeName string, headers map[string]interface{}) bool {
	err := m.mainChannel.QueueBind(
		queueName,     // queue name
		"",            // routing key
		exchangeName,  // exchange name
		false,        //no-wait
		headers)           //arguments

	if err != nil {
		log.Println("RabbitMQManager BindQueue(): Error bind queue:to exchange", err)
		return false
	}
	return true
}


func (m *RabbitMQManager) Consume(queueName string) (*amqp.Channel, <-chan amqp.Delivery, error) {
	ch, err := m.Conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)

	if err != nil {
		return nil, nil, err
	}

	messages, err := ch.Consume(
		queueName, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	if err != nil {
		return nil, nil, err
	}

	return ch, messages, nil
}

//mandatory: This flag tells the server how to react if the message cannot be routed to a queue.
// If this flag is set, the server will return an unroutable message with a Return method.
// If this flag is zero, the server silently drops the message.
//
//immediate: This flag tells the server how to react if the message cannot be routed to a queue consumer immediately.
// If this flag is set, the server will return an undeliverable message with a Return method.
// If this flag is zero, the server will queue the message, but with no guarantee that it will ever be consumed.
//The server SHOULD implement the immediate flag.
func (m *RabbitMQManager) TopicPublishText(exchangeName string, routingKey string, content []byte) bool {
	ch, err := m.Conn.Channel()
	if err != nil {
		return false
	}
	defer ch.Close()

	err = ch.Publish(
		exchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        content,
		})

	if err != nil {
		log.Println("RabbitMQManager TopicPublish(): Error publish topic message:", err)
		return false
	}

	return true
}

func (m *RabbitMQManager) TopicPublishJson(exchangeName string, routingKey string, content []byte) bool {
	ch, err := m.Conn.Channel()
	if err != nil {
		return false
	}
	defer ch.Close()

	err = ch.Publish(
		exchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        content,
		})

	if err != nil {
		log.Println("RabbitMQManager TopicPublish(): Error publish topic message:", err)
		return false
	}
	return true
}

func (m *RabbitMQManager) PublishJsonWithHeaders(exchangeName string, headers map[string]interface{}, content []byte) bool {
	ch, err := m.Conn.Channel()
	if err != nil {
		return false
	}
	defer ch.Close()

	err = ch.Publish(
		exchangeName,
		"",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers: headers,
			ContentType: "application/json",
			Body:        content,
		})

	if err != nil {
		log.Println("RabbitMQManager TopicPublish(): Error publish topic message:", err)
		return false
	}
	return true
}


func (m *RabbitMQManager) RPCPublishJson(exchangeName string, routingKey string, replyQueueName string, corrId string, content []byte) bool {
	ch, err := m.Conn.Channel()
	if err != nil {
		return false
	}
	defer ch.Close()

	err = ch.Publish(
		exchangeName,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			ReplyTo: replyQueueName,
			CorrelationId: corrId,
			Body:        content,
		})

	if err != nil {
		log.Println("RabbitMQManager RPCPublishJson(): Error publish json message:", err)
		return false
	}
	return true
}

func (m *RabbitMQManager) RPCPublishJsonWithHeaders(exchangeName string, headers map[string]interface{}, replyQueueName string, corrId string, content []byte) bool {
	ch, err := m.Conn.Channel()
	if err != nil {
		return false
	}
	defer ch.Close()

	err = ch.Publish(
		exchangeName,
		"",
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			Headers: headers,
			ContentType: "application/json",
			ReplyTo: replyQueueName,
			CorrelationId: corrId,
			Body:        content,
		})

	if err != nil {
		log.Println("RabbitMQManager RPCPublishJsonWithHeaders(): Error publish json message:", err)
		return false
	}
	return true
}

func (m *RabbitMQManager) RPCReplyJson(replyQueueName string, corrId string, content []byte) bool {
	ch, err := m.Conn.Channel()
	if err != nil {
		return false
	}
	defer ch.Close()

	err = ch.Publish(
		"",
		replyQueueName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			CorrelationId: corrId,
			Body:        content,
		})

	if err != nil {
		log.Println("RabbitMQManager RPCReplyJson(): Error reply json message:", err)
		return false
	}
	return true
}

func (m *RabbitMQManager) InspectQueue(qName string) (noAckMsgs int, consumers int, ok bool) {
	q, err := m.mainChannel.QueueInspect(qName)
	if err != nil {
		log.Println("RabbitMQManager InspectQueue(): Error inspect queue:", err)
		return 0,0, false
	}
	return q.Messages, q.Consumers, true
}


//close rabbitmq channel and connections
func (m *RabbitMQManager) Close() {
	if m.mainChannel != nil {
		err := m.mainChannel.Close()
		if err != nil {
			log.Println("RabbitMQManager Close(): Error close RabbitMQ channel:", err)
		}
	}

	if m.Conn != nil {
		m.Conn.Close()
	}
}
