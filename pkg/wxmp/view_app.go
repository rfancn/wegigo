package wxmp

import (
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
)

func (srv *WxmpServer) ViewAppIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("app index")

	context := make(map[string]interface{})
	context["enabledApps"] = srv.appManager.GetEnabledApps()
	context["appInfoMap"] = srv.appManager.GetAppInfoMap()

	srv.RespRender(w, "app.html", context)
}


func (srv *WxmpServer) ViewAppConfig(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("config app:", params.ByName("uuid"))
}

func (srv *WxmpServer) ViewAppToggle(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	appUuid :=  params.ByName("uuid")

	appName := r.FormValue("name")
	appEnabled := r.FormValue("enabled")

	log.Printf("toggle app: %s to: %s", appUuid,  appEnabled)

	toggleOk := false
	if appEnabled == "true" {
		toggleOk = srv.appManager.EnableApp(appUuid, appName)
		//do something for this app when enable app
	}else{
		toggleOk = srv.appManager.DisableApp(appUuid)
		//do something for this app when disable app
	}

	if toggleOk {
		srv.RespText(w, "success")
	}
}