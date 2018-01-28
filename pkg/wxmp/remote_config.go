package wxmp

/**
import (
	"log"
	"encoding/json"
)

import (
	"log"
	"encoding/json"
	"github.com/rfancn/wegigo/sdk/app"
	"github.com/rfancn/wegigo/sdk/etcd"
)

type WxmpRemoteConfig struct {
	//key: app uuid, value: AppConfig
	appConfigMap map[string]*app.AppConfig
	//key: app uuid, value: app name
	enabledApps  map[string]string
}


func NewRemoteConfig() *WxmpRemoteConfig {
	c := &WxmpRemoteConfig{}

	c.appConfigMap = make(map[string]*app.AppConfig)
	c.enabledApps = make(map[string]string)

	return c
}


//RefreshAppConfigMap: based on opType, update appConfigMap by key/value
func (c *WxmpRemoteConfig) RefreshAppConfigMap(opType string, key []byte, value []byte) {
	log.Printf("Before RefreshAppConfigMap: %v\n", c.appConfigMap)

	//get app uuid
	uuid  := string(key)[app.ETCD_APP_ROOT_UUID_START_POS:]

	appConfig := &app.AppConfig{}
	if err := json.Unmarshal(value, &appConfig); err != nil {
		log.Println("Error RefreshAppConfigMap(): failed unmarshal to AppConfig")
		return
	}

	switch opType {
	case "PUT":
		_, exist := c.appConfigMap[uuid]
		if exist {
			delete(c.appConfigMap, uuid)
		}
		c.appConfigMap[uuid] = appConfig
	case "DELETE":
		delete(c.appConfigMap, uuid)
	}
	log.Printf("After RefreshAppConfigMap: %v\n", c.appConfigMap)

}

func (c *WxmpRemoteConfig) RefreshEnabledApps() {
	log.Printf("Before RefreshEnabledApps: enabledApps:%v\n", c.enabledApps)
	for uuid, appConfig := range c.appConfigMap {
		_, exist := c.enabledApps[uuid]
		//if uuid not in enabled map, but it's status is "enabled", add it
		if !exist {
			if appConfig.Status == "enabled" {
				c.enabledApps[uuid] = appConfig.Name
			}
			//if uuid in enabled map, but it's status is not "enabled", remove it
		}else{
			if appConfig.Status != "enabled" {
				delete(c.enabledApps, uuid)
			}
		}
	}
	log.Printf("After RefreshEnabledApps: enabledApps:%v\n", c.enabledApps)
}

func (srv *WxmpServer) InitRemoteConfig() {
	srv.remoteConfig.appConfigMap = app.GetAppConfigMap(srv.etcd)

	srv.remoteConfig.RefreshEnabledApps()

	go srv.WatchAppConfig(srv.stopChan)
}

func (srv *WxmpServer) WatchAppConfig(stopChan chan struct{}) {
	watchChan := etcd.WatchWithPrefix(srv.etcd, app.ETCD_APP_ROOT)

	WATCH_LOOP:
	for {
		select {
		case watchResp := <-watchChan:
			for _, ev := range watchResp.Events {
				log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)

				srv.remoteConfig.RefreshAppConfigMap(ev.Type.String(), ev.Kv.Key, ev.Kv.Value)

				srv.remoteConfig.RefreshEnabledApps()
			}
		case <-stopChan:
			log.Println("Quit watch apps routine")
			break WATCH_LOOP
		}
	}
}


func (srv *WxmpServer) UpdateAppStatus(uuid string, status string) bool {
	log.Printf("Before ToggleApp: %v\n", srv.remoteConfig.enabledApps)

	ret := app.SetAppConfigField(srv.etcd, uuid, "Status", status)

	log.Printf("After ToggleApp: %v\n", srv.remoteConfig.enabledApps)
	return ret
}
**/