package wxmp

import (
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
)

func (srv *WxmpServer) ViewAdminIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("admin index")

	context := make(map[string]interface{})
	context["appInfos"] = srv.appInfos
	context["enabledUuids"] = r.Context().Value("enabledUuids").([]string)

	srv.RespRender(w, "index.html", context)
}

func (srv *WxmpServer) ViewAppAdminIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("app index")

	context := make(map[string]interface{})
	context["appInfos"] = srv.appInfos
	context["enabledUuids"] = r.Context().Value("enabledUuids").([]string)

	srv.RespRender(w, "app.html", context)
}


func (srv *WxmpServer) ViewAppConfig(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value("params").(httprouter.Params)
	log.Println("config app:", params.ByName("Uuid"))
}

func (srv *WxmpServer) ViewAppToggle(w http.ResponseWriter, r *http.Request) {
	params := r.Context().Value("params").(httprouter.Params)
	appUuid :=  params.ByName("Uuid")
	appInfo, ok := srv.appInfos[appUuid]
	if !ok {
		log.Println("Invalid app uuid:", appUuid)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	appEnabled := r.FormValue("enabled")

	log.Printf("toggle app[%s]: %s to: %s", appInfo.Name, appUuid, appEnabled)

	toggleOk := false
	if appEnabled == "true" {
		toggleOk = srv.appManager.EnableApp(appUuid, appInfo.Name)
		//do something for this app when enable app
	}else{
		toggleOk = srv.appManager.DisableApp(appUuid)
		//do something for this app when disable app
	}

	if toggleOk {
		srv.RespText(w, "success")
	}
}