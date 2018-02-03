package app

import (
	"encoding/json"
	"log"
	"path/filepath"
)

//GetEnabledAppInfos: get enabled apps, return is a map[string]string
//key is: app uuid, value is: app name
func (m *AppManager) GetEnabledAppInfos() map[string]string {
	bv := m.etcdManager.GetBytes(ETCD_APP_ENABLED_URL)

	enabledAppInfos := make(map[string]string)
	if err := json.Unmarshal(bv, &enabledAppInfos); err != nil {
		log.Printf("AppManager GetEnabledAppInfos(): Error unmarshal map[string]string:%v", err)
		return nil
	}

	return enabledAppInfos
}

//EnableApp: forcefully sync the [uuid]=name map to /app/enabled
//if someone change the app name, here forcefully sync without check if it exist or not
//will still sync the latest info to /app/enabled
func (m *AppManager) EnableApp(uuid string, name string) bool {
	item := make(map[string]string)
	item[uuid] = name
	return m.etcdManager.PutValue(filepath.Join(ETCD_APP_ENABLED_URL, uuid), item)
}

func (m *AppManager) DisableApp(uuid string) bool {
	return m.etcdManager.Delete(filepath.Join(ETCD_APP_ENABLED_URL, uuid))
}

func (m *AppManager) WatchEnabledApps(enabledApps map[string]string) (chan struct{}){
	stopChan := make(chan struct{})

	go func() {
		watchChan := m.etcdManager.Watch(ETCD_APP_ENABLED_URL)

	WATCH_LOOP:
		for {
			select {
			case watchResp := <-watchChan:
				for _, ev := range watchResp.Events {
					log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
					enabledApps = m.GetEnabledAppInfos()
				}
			case <-stopChan:
				log.Println("Quit WatchEnabledApps() routine")
				break WATCH_LOOP
			}
		}
	}()

	return stopChan
}

