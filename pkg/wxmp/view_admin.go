package wxmp

import (
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
	"bytes"
	"time"
	"encoding/json"
)

func (srv *WxmpServer) ViewAdminIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("admin index")

	context := make(map[string]interface{})
	context["appInfos"] = srv.appInfos
	context["enabledUuids"] = r.Context().Value("enabledUuids").([]string)

	srv.RespRenderFile(w, "index.html", context)
}

func (srv *WxmpServer) ViewAppAdminIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("app index")

	context := make(map[string]interface{})
	context["appInfos"] = srv.appInfos
	context["enabledUuids"] = r.Context().Value("enabledUuids").([]string)

	srv.RespRenderFile(w, "app_list.html", context)
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

//Read config stuff from yaml
func (srv *WxmpServer) GetConfigItems(yamlContent []byte) (html, inlineJs, externalJs string) {
	if ok := srv.y2h.ReadBytes(yamlContent); !ok{
		log.Println("Error parse yaml content")
		return "", "", ""
	}

	//get Javascript output
	var inlineJsBuffer bytes.Buffer
	var externalJsBuffer bytes.Buffer
	for jsType, jsContent := range srv.y2h.GetJavascript() {
		switch jsType {
		case "inline":
			inlineJsBuffer.WriteString(jsContent)
		case "external":
			externalJsBuffer.WriteString(jsContent)
		}
	}

	return srv.y2h.GetHtml(), inlineJsBuffer.String(), externalJsBuffer.String()
}

func (srv *WxmpServer) ViewAppConfigIndex(w http.ResponseWriter, r *http.Request) {
	appUuid := r.Context().Value("appUuid").(string)
	log.Println("App Config:", appUuid)
	theApp, exist := srv.enabledApps[appUuid]
	if !exist {
		log.Printf("Error config app:%s as it doesn't exist or enabled\n", appUuid)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		log.Println("display app config page for:", appUuid)

		yamlContent := theApp.GetConfigYaml()
		if yamlContent == nil {
			srv.RespRenderFile(w, "app_config.html", nil)
			return
		}

		//render template
		context := make(map[string]interface{})
		appInfo :=  theApp.GetAppInfo()
		appConfig, _ := json.Marshal(appInfo)
		context["appConfig"] = string(appConfig)
		context["appInfo"] =appInfo
		context["html"], context["inlineJs"], context["externalJs"] = srv.GetConfigItems(yamlContent)

		srv.RespRenderFile(w, "app_config.html", context)
	case "POST":
		log.Println("save app config for:", appUuid)
	}
}

func (srv *WxmpServer) ViewFetchAppConfig(w http.ResponseWriter, r *http.Request) {
	appUuid := r.Context().Value("appUuid").(string)
	log.Println("Fetch app config:", appUuid)
	_, exist := srv.enabledApps[appUuid]
	if !exist {
		log.Printf("Error config app:%s as it doesn't exist or enabled\n", appUuid)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	time.Sleep(2 * time.Second)

	srv.RespJson(w, srv.appInfos[appUuid])
}

func (srv *WxmpServer) ViewSaveAppConfig(w http.ResponseWriter, r *http.Request) {
	appUuid := r.Context().Value("appUuid").(string)
	log.Println("App Config:", appUuid)
	_, exist := srv.enabledApps[appUuid]
	if !exist {
		log.Printf("Error config app:%s as it doesn't exist or enabled\n", appUuid)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	time.Sleep(2 * time.Second)

	srv.RespJson(w, "456")
}

