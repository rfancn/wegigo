package wxmp

import (
	"github.com/rfancn/wegigo/sdk/app"
	"os"
	"path/filepath"
	"log"
	"plugin"
)

//config store
type AppLoader struct {
	srv *WxmpServer
	//store all initialized App instance
	appMap map[string]app.IApp
	//store all app info, key: app uuid, value: AppInfo
	appInfoMap map[string]*app.AppInfo

	enabledAppMap  map[string]app.IApp

	//key: app uuid, value: app name
	enabledAppInfos  map[string]string

}

func NewAppLoader(server *WxmpServer) *AppLoader {
	loader := &AppLoader{srv: server}

	loader.appMap = make(map[string]app.IApp)
	loader.appInfoMap = make(map[string]*app.AppInfo)
	loader.enabledAppMap = make(map[string]app.IApp)
	loader.enabledAppInfos = make(map[string]string)

	return loader
}

func (loader *AppLoader) Load() {
	//firstly, get enabled apps from etcd
	loader.enabledAppInfos = loader.srv.appManager.GetEnabledAppInfos()

	//discover/init/run enabled apps
	discoveredApps := loader.DiscoverApps(loader.srv.cmdArg.AppPluginDir)
	for _, app := range discoveredApps {
		//init app
		if err:= app.Init(loader.srv.Name, loader.srv.cmdArg.EtcdUrl, loader.srv.cmdArg.RabbitmqUrl); err != nil {
			log.Println("%s Error init: %s", err)
			continue
		}

		appInfo := app.GetAppInfo()
		loader.appMap[appInfo.Uuid] = app
		loader.appInfoMap[appInfo.Uuid] = appInfo
	}

	for uuid, _ := range loader.enabledAppInfos {
		app, ok := loader.appMap[uuid]
		if ok {
			loader.enabledAppMap[uuid] = app
			app.Run(loader.srv.cmdArg.AppConcurrency)
		}
	}
}


//list all modules(filename ended with .so) in all nested dir from current dir
func ListAllPlugins(dir string) (map[string]string, error) {
	if _, err := os.Stat(dir); err != nil {
		return nil, err
	}

	plugins := make(map[string]string)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".so" {
			plugins[info.Name()] = path
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return plugins, nil
}

//loadApps: list all app module files and load them
func (loader *AppLoader)  DiscoverApps(appPluginDir string) []app.IApp {
	appList := make([]app.IApp, 0)

	appPlugins, err := ListAllPlugins(appPluginDir)
	if err != nil {
		log.Printf("Error list app plugins under %s: %s\n", appPluginDir, err)
		//return empty list, but not nil
		return appList
	}

	for name, path := range appPlugins {
		plug, err := plugin.Open(path)
		if err != nil {
			log.Printf("Error open app plugin[%s]: %s\n", name, err.Error())
			continue
		}

		symbol, err := plug.Lookup(app.EXPORT_SYMBOL_NAME)
		if err != nil {
			log.Printf("Error find export symbol \"%s\" for %s\n", app.EXPORT_SYMBOL_NAME, name)
			continue
		}

		app, ok := symbol.(app.IApp)
		if !ok {
			log.Printf("%s does not implement IApp interface\n", name)
			continue
		}

		appList = append(appList, app)
	}

	return appList
}



