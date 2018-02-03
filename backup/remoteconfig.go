package backup

import (
	"log"
	"encoding/json"
)

//RefreshAppConfigMap: based on opType, update appInfoMap by key/value
func (c *ConfigStore) RefreshAppConfigMap(opType string, key []byte, value []byte) {
	log.Printf("Before RefreshAppConfigMap: %v\n", c.appInfoMap)

	//get app uuid
	uuid  := string(key)[app.ETCD_APP_ROOT_UUID_START_POS:]

	appConfig := &app.AppConfig{}
	if err := json.Unmarshal(value, &appConfig); err != nil {
		log.Println("Error RefreshAppConfigMap(): failed unmarshal to AppConfig")
		return
	}

	switch opType {
	case "PUT":
		_, exist := c.appInfoMap[uuid]
		if exist {
			delete(c.appInfoMap, uuid)
		}
		c.appInfoMap[uuid] = appConfig
	case "DELETE":
		delete(c.appInfoMap, uuid)
	}
	log.Printf("After RefreshAppConfigMap: %v\n", c.appInfoMap)

}



func (c *ConfigStore) RefreshEnabledApps() {
	log.Printf("Before RefreshEnabledApps: enabledApps:%v\n", c.enabledApps)
	for uuid, appConfig := range c.appInfoMap {
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

	go srv.WatchEnabledApps(srv.stopChan)
}

func (srv *WxmpServer) WatchEnabledApps(stopChan chan struct{}) {
	watchChan := srv.etcdManager.WatchWithPrefix(srv.etcd, app.ETCD_APP_ROOT)

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

