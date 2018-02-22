package app

import (

)

const ETCD_APP_ROOT_URL = "/app"

const ETCD_APP_INFO_URL =  "/app/info"
const ETCD_APP_ENABLED_URL = "/app/enabled"
const ETCD_APP_CONFIG_URL = "/app/config"

//exported symbol
const EXPORT_SYMBOL_NAME = "App"

/*
Registry: Register to available app list
Enable:  1.read app basic Info, 2. write to etcd  3. create queue
Disable: 1. update etcd app configuration 2. destroy proxy instance or remove the http Handler
Uninstall: 1. clean etcd app setting 2. remove from app server docker image
 */
type IApp interface {
	//serverName: the app plugin running on what kind of server
	//rootDir: app plugin .so file dir
	//etcdUrl: Etcd connection url
	//amqpUrl: RabbitMQ connection url
	//Init(env *AppEnv) error
	Init(serverName string, rootDir string, etcdUrl string, amqpUrl string) error

	//NewAppEnv() *AppEnv

	GetAppInfo() *AppInfo

	GetConfigYaml() []byte

	//check if passed data is matched or not,
	//if matched, then message will send to app's queue
	//else, no message will be recevied
	Match(data []byte) 		bool
	//start app plugin
	Start(concurrency int)
	//stop app plugin
	Stop()
	//check app is running or not
	IsRunning()				bool
	//process data
	Process(data []byte) 	[]byte
	//get own customized urls
	GetRoutes() []*AppRoute

	LoadConfig()
}

//app running Env related
type AppEnv struct {
	ServerName string
	RootDir    string
	EtcdUrl string
	AmqpUrl string
}

//App basic Info
type AppInfo struct {
	Uuid  	string
	Name  	string
	Version string
	Author 	string
	Desc  	string
	Configurable bool
}



