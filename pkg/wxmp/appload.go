package wxmp

import (
	"os"
	"path/filepath"
	"log"
	"plugin"
	"github.com/rfancn/wegigo/sdk/app"
)

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

//loadApps: list all app module files and init them
func (srv *WxmpServer)  DiscoverApps(appPluginDir string) map[string]app.IApp {
	apps := make(map[string]app.IApp)

	appPlugins, err := ListAllPlugins(appPluginDir)
	if err != nil {
		log.Printf("Error list app plugins under %s: %s\n", appPluginDir, err)
		//return empty list, but not nil
		return apps
	}

	for name, path := range appPlugins {
		log.Println("Found App Plugin:", name)

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

		//init app
		if err:= app.Init(srv.Name, srv.cmdArg.EtcdUrl, srv.cmdArg.RabbitmqUrl); err != nil {
			log.Println("%s Error init: %s", err)
			continue
		}

		apps[app.GetAppInfo().Uuid] = app
	}

	return apps
}



func (srv *WxmpServer) LoadAndRunApps() {
	//load apps: discover/init apps
	srv.apps = srv.DiscoverApps(srv.cmdArg.AppPluginDir)
	for uuid, app := range srv.apps {
		srv.appInfos[uuid] = app.GetAppInfo()
	}

	//get enabled apps from etcd
	enabledUuids := srv.appManager.GetEnabledAppUuids()
	for _, uuid := range enabledUuids {
		app, ok := srv.apps[uuid]
		if ok {
			srv.enabledApps[uuid] = app
			app.Start(srv.cmdArg.AppConcurrency)
		}
	}
}



