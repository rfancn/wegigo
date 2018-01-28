package wxmp

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

