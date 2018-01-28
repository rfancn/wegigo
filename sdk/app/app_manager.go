package app

import (
	"path/filepath"
	"log"
	"encoding/json"
	"github.com/rfancn/wegigo/sdk/etcd"
)

type AppManager struct {
	etcdManager *etcd.EtcdManager
	AppInfoMap  map[string]*AppInfo
	EnabledApps map[string]string
}

func NewAppManager(etcdManager *etcd.EtcdManager) *AppManager {
	return &AppManager{etcdManager: etcdManager}
}

func (m *AppManager) Init() {
	m.AppInfoMap = m.GetAppInfoMap()
	m.EnabledApps = m.GetEnabledApps()
}

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

//GetEnabledApps: get enabled apps, return is a map[string]string
//key is: app uuid, value is: app name
func (m *AppManager) GetEnabledApps() map[string]string {
	bv := m.etcdManager.GetBytes(ETCD_APP_ENABLED_URL)

	enabledApps := make(map[string]string)
	if err := json.Unmarshal(bv, &enabledApps); err != nil {
		log.Printf("AppManager GetEnabledApps(): Error unmarshal map[string]string:%v", err)
		return nil
	}

	return enabledApps
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

func (m *AppManager) PutAppInfo(uuid string, appInfo *AppInfo) bool {
	return m.etcdManager.PutValue(filepath.Join(ETCD_APP_INFO_URL, uuid), appInfo)
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

/**
func SetAppConfigField(cli *clientv3.Client, uuid string, field string, value interface{}) bool {
	appConfig := GetAppConfig(cli, uuid)
	if  appConfig == nil {
		log.Printf("Error SetAppConfigField(): no such app[%s]\n", uuid)
		return false
	}

	if ! setFieldValue(appConfig, field, value) {
		log.Println("Error SetAppConfigField(): failed to set field")
		return false
	}

	return etcd.Put(cli, filepath.Join(ETCD_APP_ROOT_URL, uuid), appConfig)
}

// set app config filed with reflect method
func setFieldValue(obj interface{}, field string, value interface{}) bool {
	f := reflect.ValueOf(obj).Elem().FieldByName(field)
	if !f.IsValid() || ! f.CanSet(){
		log.Println("Error setFieldValue(): field is invalid or cannot set")
		return false
	}

	switch f.Kind() {
	case reflect.String:
		sv, ok := value.(string)
		if !ok {
			log.Println("Error setFieldValue(): value is not a string")
			return false
		}
		f.SetString(sv)
	case reflect.Int64:
		iv, ok := value.(int64)
		if !ok {
			log.Println("Error setFieldValue(): value is not a int64")
			return false
		}
		f.SetInt(iv)
	}

	return true
}
**/
