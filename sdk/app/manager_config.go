package app

import (
	"path/filepath"
	"log"
)

//GetAppInfoString: get AppInfo as bytes by Uuid
func (m *AppManager) GetAppConfigBytes(uuid string) []byte {
	return m.etcdManager.GetValue(filepath.Join(ETCD_APP_CONFIG_URL, uuid))
}

func (m *AppManager) PutAppConfigBytes(uuid string, configData []byte) bool {
	return m.etcdManager.PutValueBytes(filepath.Join(ETCD_APP_CONFIG_URL, uuid), configData)
}

func (m *AppManager) WatchAppConfig(stopWatch chan struct{}, callback func(string)) {
	watchRespChan := m.etcdManager.WatchWithPrefix(ETCD_APP_CONFIG_URL)

	WATCH_LOOP:
	for {
		select {
		case watchResp := <-watchRespChan:
			for _, ev := range watchResp.Events {
				key := string(ev.Kv.Key)
				uuid := key[len(ETCD_APP_CONFIG_URL)+1:]
				callback(uuid)
			}
		case <-stopWatch:
			log.Println("Quit WatchAppConfig() routine")
			break WATCH_LOOP
		}
	}
}