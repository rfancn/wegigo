package app

import (

)

const ETCD_APP_ROOT_URL = "/app"

const ETCD_APP_INFO_URL =  "/app/info"
const ETCD_APP_ENABLED_URL = "/app/enabled"
const ETCD_APP_CONFIG_URL = "/app/config"

//uuid start position in key
const ETCD_APP_ROOT_UUID_START_POS = len(ETCD_APP_ROOT_URL) + 1




/*
Registry: Register to available app list
Enable:  1.read app basic info, 2. write to etcd  3. create queue
Disable: 1. update etcd app configuration 2. destroy proxy instance or remove the http handler
Uninstall: 1. clean etcd app setting 2. remove from app server docker image
 */
type IApp interface {
	Init(appManager *AppManager) 	error
}

//App basic info
type AppInfo struct {
	Uuid  string
	Name  string
	Version string
	Author string
	Desc  string
}

//specific app setting
type AppConfig struct {}