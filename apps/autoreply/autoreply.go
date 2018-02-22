package main

import (
	"log"
	"github.com/rfancn/wegigo/sdk/app"
	"github.com/rfancn/wegigo/sdk/wxmp"
	"time"
	"path"
	"encoding/json"
)

var App autoReplyApp
var APP_INFO = &app.AppInfo{
	Uuid: "9856d61b-e9fe-4346-b64d-dcd1102f2719",
	Name: "AutoReply",
	Version: "0.0.1",
	Author: "Ryan Fan",
	Desc: "auto reply wechat messages based on rules",
	Configurable: true,
}

type autoReplyApp struct {
	app.BaseApp
	config *AutoReplyConfig
}

type AutoReplyConfig struct {
	Keyword string
	Reply string
}

func (a *autoReplyApp) 	Init(serverName string, rootDir string, etcdUrl string, amqpUrl string) error {
	return a.BaseApp.Initialize(serverName, rootDir, etcdUrl, amqpUrl, APP_INFO, a)
}

func (a *autoReplyApp) LoadConfig() {
	configData := a.GetConfigData()
	if configData == nil {
		log.Println("Error load config data from DB!")
		return
	}

	newConfig := &AutoReplyConfig{}
	err := json.Unmarshal(configData, newConfig)
	if err != nil {
		log.Println("Error unmarshal app config data!")
		return
	}

	a.config = newConfig
}

func (a *autoReplyApp) Match(data []byte) bool{
	wxmpRequest := wxmp.NewRequest(data)
	if (wxmpRequest.MsgType == "text") && (wxmpRequest.Content == a.config.Keyword) {
		return true
	}
	return false
}

func (a *autoReplyApp) Process(data []byte) []byte{
	wxmpRequest := wxmp.NewRequest(data)
	if wxmpRequest == nil {
		log.Println("autoReplyApp Process(): Error new WxmpRequest:%s", string(data))
		return nil
	}

	log.Println("AutoReply received:", wxmpRequest.Content)

	time.Sleep(3 * time.Second)

	return wxmp.NewReply(wxmpRequest).ReplyText(a.config.Reply)
}

func (a *autoReplyApp) GetConfigYaml() []byte{
	yamlPath := path.Join(a.Env.RootDir, "asset/yaml/config.yaml")

	data, err := Asset(yamlPath)
	if err != nil {
		log.Println("AutoReplyApp: Error retrieve data:", err)
		return nil
	}

	return data
}








