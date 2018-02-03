package app

import (
	"log"
	"errors"
	"github.com/rfancn/wegigo/sdk/rabbitmq"
)

type BaseApp struct {
	serverName string
	rmqManager *rabbitmq.RabbitMQManager
	appManager *AppManager
	info       *AppInfo
	//the queue to receive message from broker
	receiveQueue string
	currentApp IApp
	rabbitMQBindHeaders map[string]interface{}
}

func (a *BaseApp)  GetAppInfo() *AppInfo {
	return a.info
}

func (a *BaseApp) Initialize(serverName string, etcdUrl string, amqpUrl string, info *AppInfo, currentApp IApp) error {
	log.Println("BaseApp Init")

	a.serverName = serverName

	rmqManager, err := rabbitmq.NewRabbitMQManager(amqpUrl)
	if err != nil {
		return err
	}

	appManager, err := NewAppManager(etcdUrl)
	if err  != nil {
		return err
	}

	a.appManager = appManager
	a.rmqManager = rmqManager
	a.currentApp = currentApp
	a.info = info
	//sync info to etcd
	if ! a.appManager.PutAppInfo(info) {
		return errors.New("Error sync AppInfo to Etcd")
	}

	//init rabbitmq bind headers which will be used later
	a.rabbitMQBindHeaders = map[string]interface{}{
		//all argument need matched
		"x-match": "all",
		//each time the server will invoke app's match function before sending message,
		//later we will set header with following key/values, and it will bind to exchange
		//1. app's uuid <-> app's match() must be true
		a.info.Uuid: "true",
	}

	return nil
}

func (a *BaseApp) Consume() {
	//1. declare exchange
	a.rmqManager.DeclareHeadersExchange(a.serverName)

	//2. declare a receive queue
	qName := a.rmqManager.DeclareTempQueue()
	if qName == "" {
		log.Fatalf("[ERROR] BaseApp Consume(): Error declare temp receive queue")
	}

	if !  a.rmqManager.BindQueueWithHeaders(qName, a.serverName, a.rabbitMQBindHeaders) {
		log.Fatalf("[ERROR] BaseApp Consume(): Error bind queue with headers")
	}

	//receive queue created by server, use the app's uuid as queue name
	ch, messages, err := a.rmqManager.Consume(qName)
	defer ch.Close()
	if err != nil {
		log.Fatalf("[ERROR] BaseApp Consume(): Error consume queue:%s\n",err)
	}

	for msg := range messages {
		//spaw a go routine to procedd message
		go a.currentApp.Process(msg.ReplyTo, msg.CorrelationId, msg.Body)
	}
}

func (a *BaseApp) Run(concurrency int) {
	log.Println("Run app:", a.info.Name)
	for j := 0; j < concurrency; j++ {
		go a.Consume()
	}
}

func (a BaseApp) Response(replyQueueName string, corrId string, data []byte) {
	//a.rmqManager.TopicPublishJson(a.serverName, "reply."+ a.info.Uuid, data)
	a.rmqManager.RPCReplyJson(replyQueueName, corrId, data)
}
