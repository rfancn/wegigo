package main

import (
	"log"
	"github.com/rfancn/wegigo/sdk/app"
	"github.com/rfancn/wegigo/sdk/wxmp"
	"time"
)

//go:generate go-bindata -o apps/autoreply/bindata.go apps/autoreply/asset/...
var App publishToHugo
var APP_INFO = &app.AppInfo{
	Uuid: "a419fe78-6c0a-4c97-a6d7-735cf03cfe2d",
	Name: "HugoPublish",
	Version: "0.0.1",
	Author: "Ryan Fan",
	Desc: "Automatically publish url to Hugo site based on received Wechat messages",
	Configurable: true,
}

type publishToHugo struct {
	app.BaseApp
}

func (a *publishToHugo) Init(serverName string, rootDir string, etcdUrl string, amqpUrl string) error {
	return a.BaseApp.Initialize(serverName, rootDir, etcdUrl, amqpUrl, APP_INFO, a)
}

func (a *publishToHugo) Match(data []byte) bool{
	return true
}

func (a *publishToHugo) Process(data []byte) []byte{
	wxmpRequest := wxmp.NewRequest(data)
	if wxmpRequest == nil {
		log.Println("autoReplyApp Process(): Error new WxmpRequest:%s", string(data))
		return nil
	}

	log.Println("HugoAutoPublish received:", wxmpRequest.Content)

	time.Sleep(1 * time.Second)

	return wxmp.NewReply(wxmpRequest).ReplyText("hugo:" + wxmpRequest.Content)
}







