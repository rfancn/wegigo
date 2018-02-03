package main

import (
	"log"
	"github.com/rfancn/wegigo/sdk/app"
	"encoding/json"
	"time"
)

//go:generate go-bindata -o apps/autoreply/bindata.go apps/autoreply/asset/...

var APP_INFO = &app.AppInfo{
	Uuid: "a419fe78-6c0a-4c97-a6d7-735cf03cfe2d",
	Name: "HugoAutoPublish",
	Version: "0.0.1",
	Author: "Ryan Fan",
	Desc: "Automatically publish page to Hugo site based on received Wechat messages",
}

type publishToHugo struct {
	app.BaseApp
}

func (a *publishToHugo) Init(appManager *app.AppManager) error {
	log.Println("hugoAutoPublish Init")

	a.BaseApp.Initialize(appManager, APP_INFO)

	return nil
}

func (a *publishToHugo) Proceed(data []byte) []byte{
	wxmpRequest := &wxmp.Request{}
	err := json.Unmarshal(data, &wxmpRequest)
	if err != nil {
		log.Println("Error unmarsh amqp message to WxmpRequest:", err)
		return nil
	}

	log.Println("AutoReply received:", wxmpRequest.Content)

	time.Sleep(3 * time.Second)

	a.Response(wxmp.NewReply(wxmpRequest).ReplyText("echo:" + wxmpRequest.Content))

	return nil
}

var App publishToHugo





