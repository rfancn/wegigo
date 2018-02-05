package demo

import (
	"net/http"
	//"bytes"
	"log"
	"github.com/julienschmidt/httprouter"
	"context"

)

//convert http.handler to httprouter.handle
func HandleFunc(h http.Handler, enabledApps map[string]string) httprouter.Handle {
	log.Println("in handleFunc")
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "enabledApps", enabledApps)
		ctx = context.WithValue(ctx, "params", params)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}

func isEnabledMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("is enabled handler invoked")
		enabledApps, ok := r.Context().Value("enabledApps").(map[string]string)
		if !ok {
			log.Println("enabledApps is not a map[string]string")
			return
		}

		params, ok := r.Context().Value("params").(httprouter.Params)
		if !ok {
			log.Println("Params is not a httprouter.Params")
			return
		}

		reqAppId := params.ByName("Uuid")

		log.Println(reqAppId)
		log.Println(enabledApps)

		_, enabled := enabledApps[reqAppId]
		if !enabled {
			log.Println("App is not enabled, skip...")
			return
		}

		next.ServeHTTP(w, r)
	})
}


func mainHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("process index")
}

