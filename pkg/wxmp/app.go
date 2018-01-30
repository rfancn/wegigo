package wxmp

import (
	"os"
	"path/filepath"
	"log"
	"plugin"
	AppApi "github.com/rfancn/wegigo/sdk/app"
)

//list all modules(filename ended with .so) in all nested dir from current dir
func listAllModules(dir string) map[string]string {
	modules := make(map[string]string)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".so" {
			modules[info.Name()] = path
		}

		return nil
	})

	if err != nil {
		log.Printf("Error listAllFiles(): [%v]\n", err)
	}

	return modules
}

//loadApps: list all app module files and load them
func (srv *WxmpServer) loadApps() {
	if _, err := os.Stat(srv.appsDir); err != nil {
		log.Fatal("Non exist apps dir:", srv.appsDir)
	}


	appModuleInfos := listAllModules(srv.appsDir)
	for appModName, appModPath := range appModuleInfos {
		plug, err := plugin.Open(appModPath)
		if err != nil {
			log.Printf("failed to open app module[%s]:%s", appModName, err.Error())
			continue
		}

		symbol, err := plug.Lookup("App")
		if err != nil {
			log.Printf("%s does not export symbol \"%s\"\n", appModName, "App")
			continue
		}

		app, ok := symbol.(AppApi.IApp)
		if !ok {
			log.Printf("%s does not implement IApp interface\n", appModName)
			continue
		}

		if err := app.Init(srv.appManager); err != nil {
			log.Printf("%s initialization failed: %s\n", appModName, err.Error())
			continue
		}

		//run app
		go app.Run()
	}


	//srv.ctx = context.WithValue(srv.ctx, "apps", apps)
}

func (srv *WxmpServer) discoverApps() {
	srv.rmqManager.DeclareTopicExchange(SERVER_NAME, false)

	srv.loadApps()
	//srv.appManager.Init()
}
