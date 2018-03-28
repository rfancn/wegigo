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
	if bv == nil {
		log.Println("AppManager GetEnabledAppUuids(): Empty enabled apps")
		return nil
	}

	enabledAppKVs := make(map[string]string)
	if err := json.Unmarshal(bv, &enabledAppKVs); err != nil {
		log.Printf("AppManager GetEnabledAppUuids(): Error unmarshal map[string]string:%v", err)
		return nil
	}

	return enabledAppKVs
}

//GetEnabledAppUuids: get enabled app Uuids
func (m *AppManager) GetEnabledAppUuids() []string {
	//create new empty enabledAppIds
	enabledAppIds := make([]string,0)

	enabledAppKVs := m.GetEnabledAppKVs()
	// if we failed to get enabled App key/values, return empty enabledAppIds
	if enabledAppKVs == nil {
		return enabledAppIds
	}

	for Uuid, _ := range enabledAppKVs {
		enabledAppIds = append(enabledAppIds, Uuid)
	}

	return enabledAppIds
}

//EnableApp: forcefully sync the [Uuid]=name map to /app/enabled
//if someone change the app name, here forcefully sync without check if it exist or not
//will still sync the latest Info to /app/enabled
func (m *AppManager) EnableApp(Uuid string, name string) bool {
	kvs := m.GetEnabledAppKVs()
	kvs[Uuid] = name
	return m.etcdManager.TxnPutValueAny(filepath.Join(ETCD_APP_ENABLED_URL), kvs)
}

func (m *AppManager) DisableApp(Uuid string) bool {
	kvs := m.GetEnabledAppKVs()

	_, ok := kvs[Uuid];
	if ok {
		delete(kvs, Uuid);
	}

	return m.etcdManager.TxnPutValueAny(filepath.Join(ETCD_APP_ENABLED_URL), kvs)
}

func (m *AppManager) WatchEnabledApps(stopWatch chan struct{}, callback func()) {
	watchRespChan := m.etcdManager.Watch(ETCD_APP_ENABLED_URL)

	WATCH_LOOP:
	for {
		select {
		case watchResp := <-watchRespChan:
			for _, ev := range watchResp.Events {
				log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
				callback()
			}
		case <-stopWatch:
			log.Println("Quit WatchEnabledApps() routine")
			break WATCH_LOOP
		}
	}
}
