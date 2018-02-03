package app

import (
	"path/filepath"
	"log"
	"encoding/json"
)

//GetAppInfoString: get AppInfo as bytes by uuid
func (m *AppManager) GetAppInfoBytes(uuid string) []byte {
	return m.etcdManager.GetBytes(filepath.Join(ETCD_APP_INFO_URL, uuid))
}

//GetAppInfo: get AppInfo object from etcd by uuid
func (m *AppManager) GetAppInfo(uuid string) *AppInfo {
	appInfoBytes := m.GetAppInfoBytes(uuid)
	if appInfoBytes == nil {
		log.Println("AppManager GetAppInfo(): No such app info for:", uuid)
		return nil
	}

	//try unmarshal it
	appInfo := &AppInfo{}
	if err := json.Unmarshal(appInfoBytes, &appInfo); err != nil {
		log.Printf("AppManager GetAppInfo(): Error unmarshal AppInfo:%v", err)
		return nil
	}

	return appInfo
}


//GetAppInfoMap: get current appconfig map in etcd
// key: app uuid, value: AppConfig
func (m *AppManager) GetAppInfoMap() map[string]*AppInfo{
	appInfoList := m.etcdManager.GetBytesList(ETCD_APP_INFO_URL)

	appInfoMap := make(map[string]*AppInfo)
	for _, appInfoBytes := range appInfoList {
		appInfo := &AppInfo{}
		if err := json.Unmarshal(appInfoBytes, &appInfo); err != nil {
			log.Println("AppManager GetAppInfoMap(): Error unmarshal AppInfo:", err)
			continue
		}

		appInfoMap[appInfo.Uuid] = appInfo
	}

	return appInfoMap
}

func (m *AppManager) PutAppInfo(appInfo *AppInfo) bool {
	return m.etcdManager.PutValue(filepath.Join(ETCD_APP_INFO_URL, appInfo.Uuid), appInfo)
}
