package app

import (
	"path/filepath"
	"log"
	"encoding/json"
)

//GetAppInfoString: get AppInfo as bytes by Uuid
func (m *AppManager) GetAppInfoBytes(Uuid string) []byte {
	return m.etcdManager.GetValue(filepath.Join(ETCD_APP_INFO_URL, Uuid))
}

//GetAppInfo: get AppInfo object from etcd by Uuid
func (m *AppManager) GetAppInfo(Uuid string) *AppInfo {
	appInfoBytes := m.GetAppInfoBytes(Uuid)
	if appInfoBytes == nil {
		log.Println("AppManager GetAppInfo(): No such app info for:", Uuid)
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
// key: app Uuid, value: AppConfig
func (m *AppManager) GetAppInfoMap() map[string]*AppInfo{
	appInfoList := m.etcdManager.GetValueList(ETCD_APP_INFO_URL)

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
	return m.etcdManager.PutValueAny(filepath.Join(ETCD_APP_INFO_URL, appInfo.Uuid), appInfo)
}
