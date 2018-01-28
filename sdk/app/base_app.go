package app

import (
	"encoding/json"
	"log"
)

type BaseApp struct {
	appManager *AppManager
	appInfo *AppInfo
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

	//sync with etcd
	a.SyncWithEtcd()
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
