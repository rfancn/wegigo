package main

import (
	"log"
	"github.com/rfancn/wegigo/sdk/app"
	"github.com/rfancn/wegigo/sdk/wxmp"
	"time"
)

var APP_INFO = &app.AppInfo{
	Uuid: "9856d61b-e9fe-4346-b64d-dcd1102f2719",
	Name: "AutoReply",
	Version: "0.0.1",
	Author: "Ryan Fan",
	Desc: "auto reply wechat messages based on rules",
}

type autoReplyApp struct {
	app.BaseApp
}

func (a *autoReplyApp) Init(serverName string, etcdUrl string, amqpUrl string) error {
	log.Println("autoReplyApp Init")
	return a.BaseApp.Initialize(serverName, etcdUrl, amqpUrl, APP_INFO, a)
}

func (a *autoReplyApp) Match(data []byte) bool{
	wxmpRequest := wxmp.NewRequest(data)
	if wxmpRequest.MsgType == "text" && wxmpRequest.Content == "test" {
		return true
	}
	return false
}

func (a *autoReplyApp) Process(replyQueueName string, correlationId string, data []byte) {
	wxmpRequest := wxmp.NewRequest(data)
	if wxmpRequest == nil {
		log.Println("autoReplyApp Process(): Error new WxmpRequest:%s", string(data))
		return
	}

	log.Println("AutoReply received:", wxmpRequest.Content)

	time.Sleep(3 * time.Second)

	replyContent := wxmp.NewReply(wxmpRequest).ReplyText("echo:" + wxmpRequest.Content)
	a.Response(replyQueueName, correlationId, replyContent)
}

var App autoReplyApp






