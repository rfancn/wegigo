package main

import (
	"log"
	"github.com/rfancn/wegigo/sdk/app"
)

//go:generate go-bindata -o apps/autoreply/bindata.go apps/autoreply/asset/...
//var FOR_SERVER = "wxmp"

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

func (a *autoReplyApp) Init(appManager *app.AppManager) error {
	log.Println("autoReplyApp Init")

	a.BaseApp.Initialize(appManager, APP_INFO)

	return nil
}

var App autoReplyApp







