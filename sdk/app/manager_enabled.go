package app

import (
	"encoding/json"
	"log"
	"path/filepath"
)

//GetEnabledAppKVs: get enabled apps, return is a map[string]string
//key is: app Uuid, value is: app name
func (m *AppManager) GetEnabledAppKVs() map[string]string {
	bv := m.etcdManager.GetValue(ETCD_APP_ENABLED_URL)

	enabledAppKVs := make(map[string]string)
	if err := json.Unmarshal(bv, &enabledAppKVs); err != nil {
		log.Printf("AppManager GetEnabledAppUuids(): Error unmarshal map[string]string:%v", err)
		return enabledAppKVs
	}

	return enabledAppKVs
}

//GetEnabledAppUuids: get enabled app Uuids
func (m *AppManager) GetEnabledAppUuids() []string {
	enabledAppKVs := m.GetEnabledAppKVs()

	enabledAppIds := make([]string,0)
	for Uuid, _ := range enabledAppKVs {
		enabledAppIds = append(enabledAppIds, Uuid)
	}

	return enabledAppIds
}

//EnableApp: forcefully sync the [Uuid]=name map to /app/enabled
//if someone change the app name, here forcefully sync without check if it exist or not
//will still sync the latest info to /app/enabled
func (m *AppManager) EnableApp(Uuid string, name string) bool {
	kvs := m.GetEnabledAppKVs()
	kvs[Uuid] = name
	return m.etcdManager.PutValue(filepath.Join(ETCD_APP_ENABLED_URL), kvs)
}

func (m *AppManager) DisableApp(Uuid string) bool {
	kvs := m.GetEnabledAppKVs()

	_, ok := kvs[Uuid];
	if ok {
		delete(kvs, Uuid);
	}

	return m.etcdManager.PutValue(filepath.Join(ETCD_APP_ENABLED_URL), kvs)
}

/**
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
					enabledApps = m.GetEnabledAppUuids()
				}
			case <-stopChan:
				log.Println("Quit WatchEnabledApps() routine")
				break WATCH_LOOP
			}
		}
	}()

	return stopChan
}
**/
