package wxmp

import (
	"net/http"
	//"bytes"
	"log"
	"context"
	"github.com/rfancn/wegigo/sdk/wxmp"
	"github.com/justinas/alice"
	"io/ioutil"
	"strconv"
)

const TOKEN = "laonabuzhai"

//WxmpRequestMiddleware: convert http.Request to WxmpRequest
func CheckSignatureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		timestamp := r.FormValue("timestamp")
		nonce := r.FormValue("nonce")
		signature := r.FormValue("signature")

		if !wxmp.CheckSignature(TOKEN, timestamp, nonce, signature) {
			return
		}

		next.ServeHTTP(w, r)
	})
}

//WxmpRequestMiddleware: convert http.Request to WxmpRequest
func (srv *WxmpServer) MsgHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get http body
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error read http body ", err)
			w.Write([]byte("success"))
			return
		}

		msgHeaders := make(map[string]interface{})
		for uuid, app := range srv.appLoader.enabledAppMap {
			msgHeaders[uuid] = strconv.FormatBool(app.Match(data))
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "data", data)
		ctx = context.WithValue(ctx, "msgHeaders", msgHeaders)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (srv *WxmpServer) setupRouter() {
	log.Println("Setup router")

	//handle wxmp server verification request
	srv.AddHttpHandler("get", "/",
		alice.New(CheckSignatureMiddleware).Then(http.HandlerFunc(srv.ViewVerifyWxmpServer)))

	//handle wxmp server post request, request follow below logic flow:
	//1. WxmpRequestMiddleware:  HttpRequest -> WxmpRequest
	srv.AddHttpHandler("post", "/", alice.New(
		CheckSignatureMiddleware,
		srv.MsgHeaderMiddleware).Then(http.HandlerFunc(srv.ViewMain)))


	//admin urls
	srv.AddHttpHandlerFunc("get", "/admin/", srv.ViewAdminIndex)
	srv.AddHttpHandlerFunc("get", "/admin/app/", srv.ViewAppAdminIndex)

	//app config route
	srv.AddHttpHandlerFunc("get", "/admin/app/config/:uuid", srv.ViewAppConfig)
	srv.AddHttpHandlerFunc("post", "/admin/app/toggle/:uuid", srv.ViewAppToggle)
}
