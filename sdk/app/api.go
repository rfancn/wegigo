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
Enable:  1.read app basic info, 2. write to etcd  3. create queue
Disable: 1. update etcd app configuration 2. destroy proxy instance or remove the http handler
Uninstall: 1. clean etcd app setting 2. remove from app server docker image
 */
type IApp interface {
	Init(serverName string, etcdUrl string, amqpUrl string) error
	GetAppInfo() *AppInfo
	//check if passed data is matched or not,
	//if matched, then message will send to app's queue
	//else, no message will be recevied
	Match(data []byte) 		bool
	Start(concurrency int)
	Stop()
	Process(data []byte) 	[]byte
}

//App basic info
type AppInfo struct {
	Uuid  	string
	Name  	string
	Version string
	Author 	string
	Desc  	string
}



