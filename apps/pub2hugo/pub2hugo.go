package main

import (
	"log"
	"github.com/rfancn/wegigo/sdk/app"
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

var App publishToHugo





