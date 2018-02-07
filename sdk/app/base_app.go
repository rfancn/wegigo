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
	//stop channel to indicate app stop
	stopChannel chan int

	//running status
	status string
}

func (a *BaseApp)  GetAppInfo() *AppInfo {
	return a.info
}

func (a *BaseApp) Initialize(serverName string, etcdUrl string, amqpUrl string, info *AppInfo, currentApp IApp) error {
	rmqManager, err := rabbitmq.NewRabbitMQManager(amqpUrl)
	if err != nil {
		return err
	}

	appManager, err := NewAppManager(etcdUrl)
	if err  != nil {
		return err
	}

	a.serverName = serverName
	a.appManager = appManager
	a.rmqManager = rmqManager
	a.currentApp = currentApp
	a.info = info
	a.stopChannel = make(chan int)

	//sync info to etcd
	if ! a.appManager.PutAppInfo(info) {
		return errors.New("Error sync AppInfo to Etcd")
	}

	return nil
}

func (a *BaseApp) Consume(headers map[string]interface{}) {
	//1. declare a receive queue
	qName := a.rmqManager.DeclareTempQueue()
	if qName == "" {
		log.Fatalf("[ERROR] BaseApp Consume(): Error declare temp receive queue")
	}

	//2. bind headers to queue
	if !  a.rmqManager.BindQueueWithHeaders(qName, a.serverName, headers) {
		log.Fatalf("[ERROR] BaseApp Consume(): Error bind queue with headers")
	}

	//3. consume queue
	ch, messages, err := a.rmqManager.Consume(qName)
	defer ch.Close()
	if err != nil {
		log.Fatalf("[ERROR] BaseApp Consume(): Error consume queue:%s\n",err)
	}

	//4. monitor for received messages
	RECV_LOOP:
	for {
		select {
		case msg := <- messages:
			//spaw a go routine to procedd message
			go func() {
				replyData := a.currentApp.Process(msg.Body)
				if replyData != nil {
					a.rmqManager.RPCReplyJson(msg.ReplyTo, msg.CorrelationId, replyData)
				}else{
					// if Process() returns nil, then return empty []byte
					a.rmqManager.RPCReplyJson(msg.ReplyTo, msg.CorrelationId, []byte(""))
				}
			}()
		case <-a.stopChannel:
			log.Println("Consume stop")
			break RECV_LOOP
		}
	}
}

func (a *BaseApp) Start(concurrency int) {
	log.Println("Start app:", a.info.Name)
	a.status = "running"
	//1. declare exchange
	a.rmqManager.DeclareHeadersExchange(a.serverName)


	//2. build match message headers
	//init rabbitmq bind headers which will be used later
	headers := map[string]interface{}{
		//all argument need matched
		"x-match": "all",
		//each time the server will invoke app's match function before sending message,
		//later we will set header with following key/values, and it will bind to exchange
		//1. app's Uuid <-> app's match() must be true
		a.info.Uuid: "true",
	}

	for j := 0; j < concurrency; j++ {
		go a.Consume(headers)
	}
}

func (a *BaseApp) Stop() {
	log.Println("Stop app:", a.info.Name)
	a.stopChannel <- 1
	a.status = "stopped"
}

func (a *BaseApp) IsRunning() bool {
	if a.status == "running" {
		return true
	}
	return false
}

