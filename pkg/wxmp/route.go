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
	"path"
)

const TOKEN = "laonabuzhai"

//CheckSignatureMiddleware: validate the signature is correct or not
func CheckSignatureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//r.ParseForm()
		timestamp := r.FormValue("timestamp")
		nonce := r.FormValue("nonce")
		signature := r.FormValue("signature")

		if !wxmp.CheckSignature(TOKEN, timestamp, nonce, signature) {
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "nonce", nonce)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

//MsgHeaderMiddleware: get wxmp message header
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
		for Uuid, app := range srv.enabledApps {
			msgHeaders[Uuid] = strconv.FormatBool(app.Match(data))
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "data", data)
		ctx = context.WithValue(ctx, "msgHeaders", msgHeaders)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

//Default Admin middleware: pass appInfos and enabledAppUuids
func (srv *WxmpServer) AdminDefaultMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enabledAppUuids := make([]string, 0)
		for uuid, _ := range srv.enabledApps {
			enabledAppUuids = append(enabledAppUuids, uuid)
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "enabledUuids", enabledAppUuids)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

//Default Admin middleware: pass appInfos and enabledAppUuids
func (srv *WxmpServer) AppDefaultMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//pass appPluginDir as AppRoot, which will used to load go-bindata for app
		ctx := r.Context()
		ctx = context.WithValue(ctx, "AppRoot", srv.cmdArg.AppPluginDir)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}


func (srv *WxmpServer) SetupAppRoutes() {
	for uuid, app := range srv.apps {
		appRoutes := app.GetRoutes()
		for _, route := range appRoutes {
			appPrefixUrl := path.Join("/", "app", uuid, route.Url)
			srv.AddHttpHandler(
				route.Method,
				//add "/app/uuid/" as a prefix for the app route url
				appPrefixUrl,
				alice.New(
					srv.AppDefaultMiddleware,
				).Then(http.HandlerFunc(route.Handler)),
			)
		}
	}
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
	srv.AddHttpHandler("get", "/admin/",	alice.New(
		srv.AdminDefaultMiddleware).Then(http.HandlerFunc(srv.ViewAdminIndex)))

	srv.AddHttpHandler("get", "/admin/app/",	alice.New(
		srv.AdminDefaultMiddleware).Then(http.HandlerFunc(srv.ViewAppAdminIndex)))

	//app enable/disable route
	srv.AddHttpHandler("post", "/admin/app/toggle/:Uuid",	alice.New(
		srv.AdminDefaultMiddleware).Then(http.HandlerFunc(srv.ViewAppToggle)))

	srv.SetupAppRoutes()
}
