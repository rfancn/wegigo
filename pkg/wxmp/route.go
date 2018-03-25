package wxmp

import (
	"net/http"
	//"bytes"
	"log"
	"context"
	"github.com/rfancn/wegigo/sdk/wxmp"
	"github.com/justinas/alice"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"strconv"
	"path"
	"strings"
)

const TOKEN = "laonabuzhai"

//CheckSignatureMiddleware: validate the signature is correct or not
func CheckSignatureMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Check signature")

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
		//make sure we don't respond to disabled app admin requests
		params := r.Context().Value("params").(httprouter.Params)
		appUuid :=  params.ByName("Uuid")
		if appUuid != "" {
			_, exist := srv.enabledApps[appUuid]
			if !exist {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
		}

		//get enabled app uuids
		enabledAppUuids := make([]string, 0)
		for uuid, _ := range srv.enabledApps {
			enabledAppUuids = append(enabledAppUuids, uuid)
		}

		ctx := r.Context()
		//store current requested app uuid
		ctx = context.WithValue(ctx, "appUuid", appUuid)
		//store all enabled uuids
		ctx = context.WithValue(ctx, "enabledUuids", enabledAppUuids)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

//Default App middleware: pass appInfos and enabledAppUuids
func (srv *WxmpServer) AppDefaultMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enabledUrl := false
		for uuid, _ := range srv.enabledApps {
			if strings.Contains(r.URL.Path, uuid) {
				enabledUrl = true
				break
			}
		}

		//on enabled url can continue executed by app plugin,
		//otherwise, just returns
		if ! enabledUrl {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, r)
	})
}


func (srv *WxmpServer) SetupAppRoutes() {
	for uuid, app := range srv.apps {
		appRoutes := app.GetRoutes()
		for _, route := range appRoutes {
			appPrefixUrl := path.Join("/", "app", uuid, route.Url)
			srv.AddServerHandler(
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

func (srv *WxmpServer) SetupAdminRoutes() {
	//admin urls
	srv.AddServerHandler("get", "/admin/",	alice.New(
		srv.AdminDefaultMiddleware).Then(http.HandlerFunc(srv.ViewAdminIndex)))

	srv.AddServerHandler("get", "/admin/app/",	alice.New(
		srv.AdminDefaultMiddleware).Then(http.HandlerFunc(srv.ViewAppAdminIndex)))

	//app enable/disable route
	srv.AddServerHandler("post", "/admin/app/toggle/:Uuid",	alice.New(
		srv.AdminDefaultMiddleware).Then(http.HandlerFunc(srv.ViewAppToggle)))

	//add app's config url
	//app config index url
	srv.AddServerHandler("get", "/admin/app/config/:Uuid",	alice.New(
		srv.AdminDefaultMiddleware).Then(http.HandlerFunc(srv.ViewAppConfigIndex)))

	//fetch app config from db
	srv.AddServerHandler("get", "/app/config/:Uuid",	alice.New(
		srv.AdminDefaultMiddleware).Then(http.HandlerFunc(srv.ViewFetchAppConfig)))
	//save app config to db
	srv.AddServerHandler("post", "/app/config/:Uuid",	alice.New(
		srv.AdminDefaultMiddleware).Then(http.HandlerFunc(srv.ViewSaveAppConfig)))
}

func (srv *WxmpServer) SetupRouter() {
	log.Println("Setup router")

	srv.AddHandlerFunc("get", "/", srv.ViewVerifyWxmpServer)

	//handle wxmp server verification request
	srv.AddServerHandler("get", "/",
		alice.New(CheckSignatureMiddleware).Then(http.HandlerFunc(srv.ViewVerifyWxmpServer)))

	//handle wxmp server post request, request follow below logic flow:
	//1. WxmpRequestMiddleware:  HttpRequest -> WxmpRequest
	srv.AddServerHandler("post", "/", alice.New(
		CheckSignatureMiddleware,
		srv.MsgHeaderMiddleware).Then(http.HandlerFunc(srv.ViewMain)))

	srv.SetupAdminRoutes()

	srv.SetupAppRoutes()
}
