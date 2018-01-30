package app

import (
	"encoding/json"
	"log"
	"github.com/rfancn/wegigo/sdk/rabbitmq"
)

type BaseApp struct {
	rmqManager *rabbitmq.RabbitMQManager
	appManager *AppManager
	appInfo *AppInfo
	queueName string
}


func (a *BaseApp) Initialize(appManager *AppManager, appInfo *AppInfo) {
	log.Println("BaseApp Init")

	if appManager == nil {
		log.Fatal("BaseApp Initialize(): Null app manager")
	}

	if appInfo == nil {
		log.Fatal("BaseApp Initialize(): Null app info")
	}

	a.appManager = appManager
	a.appInfo = appInfo

	a.SetupRabbitMQ()

	//sync with etcd
	a.SyncWithEtcd()
}

func (a *BaseApp) SetupRabbitMQ() {
	a.rmqManager = rabbitmq.NewRabbitMQManager("127.0.0.1", 5672)

	//1. declare exchange, as we already done in server initializing step, skip it here
	a.rmqManager.DeclareTopicExchange("wxmp", false)

	//2. desclar queue
	a.queueName = a.rmqManager.DeclareQueue(a.appInfo.Uuid, false)
	if a.queueName == "" {
		log.Fatal("BaseApp SetupRabbitMQ(): Error DeclareQueue")
	}

	//3. bind queue to exchange
	if ! a.rmqManager.BindQueue(a.queueName, "wxmp", a.appInfo.Name) {
		log.Fatal("BaseApp SetupRabbitMQ(): Error BindQueue")
	}
}

func (a BaseApp) Run() {
	log.Println("Run app:", a.appInfo.Name)

	msgs := a.rmqManager.Consume(a.queueName)
	for m := range msgs {
		log.Printf(" [x] %s", m.Body)
	}
}

//GetAppInfoBytes: json marshal AppInfo and returns the string
func (a *BaseApp) GetAppInfoBytes() []byte {
	bv,err := json.Marshal(a.appInfo)
	if err != nil {
		log.Printf("BaseApp GetAppInfoBytes(): Error marshal AppInfo:%v", err)
		return nil
	}

	return bv
}

func (a *BaseApp) SyncWithEtcd() {
	dbAppInfoBytes := a.appManager.GetAppInfoBytes(a.appInfo.Uuid)
	if dbAppInfoBytes == nil {
		a.appManager.PutAppInfo(a.appInfo.Uuid, a.appInfo)
		return
	}

	//if appinfo fetched from db don't equal to the current one
	//sync the curernt one to db
	if string(dbAppInfoBytes) != string(a.GetAppInfoBytes()) {
		a.appManager.PutAppInfo(a.appInfo.Uuid, a.appInfo)
	}
}
