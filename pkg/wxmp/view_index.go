package wxmp

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"log"
)

func (srv *WxmpServer) ViewIndex(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	log.Println("server index")

	context := make(map[string]interface{})
	context["enabledApps"] = srv.appManager.GetEnabledApps()

	srv.RespRender(w, "index.html", context)
}
