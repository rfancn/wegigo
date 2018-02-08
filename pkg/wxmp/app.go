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
		log.Printf("Found App Plugin:%s under:%s", name, path)

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
		if err:= app.Init(srv.Name, filepath.Dir(path), srv.cmdArg.EtcdUrl, srv.cmdArg.RabbitmqUrl); err != nil {
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

func inSlice(item string, list []string) bool  {
	for _, i := range list {
		if item == i {
			return true
		}
	}

	return false
}


func (srv *WxmpServer) UpdateEnabledApps() {
	log.Println("Previous Enabled Apps:", srv.enabledApps)

	//get enabled apps from etcd
	newEnabledUuids := srv.appManager.GetEnabledAppUuids()

	//check need to be removed/stopped one
	for oldUuid, app := range srv.enabledApps {
		//if old uuid not in current enabled Uuid, then it need stop
		if !inSlice(oldUuid, newEnabledUuids) {
			log.Println("stop uuid:", oldUuid)
			delete(srv.enabledApps, oldUuid)
			app.Stop()
		}
	}

	//check new added/run one
	for _, newUuid := range newEnabledUuids {
		theApp, ok := srv.enabledApps[newUuid]
		//in case enabled uuid is not in current enabledApps
		if !ok {

			theApp = srv.apps[newUuid]
			srv.enabledApps[newUuid] = theApp
		}

		//start the new
		if ! theApp.IsRunning() {
			log.Println("start uuid:", newUuid)
			theApp.Start(srv.cmdArg.AppConcurrency)
		}
	}
}

func (srv *WxmpServer) WatchEnabledApps() {
	srv.appManager.WatchEnabledApps(srv.stopWatch, srv.UpdateEnabledApps)
}


