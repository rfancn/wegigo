package app

import (
	"github.com/rfancn/wegigo/sdk/etcd"
)

//AppHelper mostly used by app
type AppManager struct {
	etcdManager *etcd.EtcdManager
}

func NewAppManager(etcdUrl string) (*AppManager, error) {
	etcdManager, err := etcd.NewEtcdManager(etcdUrl)
	if err != nil {
		return nil, err
	}

	return &AppManager{etcdManager: etcdManager}, nil
}

func (m *AppManager) Close() {
	m.etcdManager.Close()
}

/**
func SetAppConfigField(cli *clientv3.Client, Uuid string, field string, value interface{}) bool {
	appConfig := GetAppConfig(cli, Uuid)
	if  appConfig == nil {
		log.Printf("Error SetAppConfigField(): no such app[%s]\n", Uuid)
		return false
	}

	if ! setFieldValue(appConfig, field, value) {
		log.Println("Error SetAppConfigField(): failed to set field")
		return false
	}

	return etcd.Put(cli, filepath.Join(ETCD_APP_ROOT_URL, Uuid), appConfig)
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
