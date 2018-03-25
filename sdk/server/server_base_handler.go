package server

import (
	"strings"
	"path/filepath"
	"log"
	"net/http"
	"context"
	"github.com/julienschmidt/httprouter"
)

//add httprouter handle func
func (srv *BaseServer) AddHttpRouterHandler(method string, withServerPrefix bool, url string, handleFunc httprouter.Handle) bool {
	lowerMethod := strings.ToLower(method)

	//with server prefix will add server name before url
	var absUrl string
	if withServerPrefix {
		absUrl = filepath.Join("/", srv.Name, url)
	}else{
		absUrl = filepath.Join("/", url)
	}

	switch lowerMethod {
	case "get":
		srv.router.GET(absUrl, handleFunc)
	case "post":
		srv.router.POST(absUrl, handleFunc)
	default:
		log.Println("Non supported http method: ", method)
		return false
	}

	return true
}

//convert http.Handler to httprouter.Handle
func NewHttpRouterHandle(h http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "params", params)
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	}
}

//standard http handler/handler func, no server name before url
func (srv *BaseServer) AddHandler(method string, url string, handler http.Handler) bool {
	return srv.AddHttpRouterHandler(method, false, url, NewHttpRouterHandle(handler))
}

func (srv *BaseServer) AddHandlerFunc(method string, url string, handlerFunc func(http.ResponseWriter, *http.Request)) bool {
	return srv.AddHttpRouterHandler(method, false, url, NewHttpRouterHandle(http.HandlerFunc(handlerFunc)))
}

//For all server handler/handler func, add server name before url
func (srv *BaseServer) AddServerHandler(method string, url string, handler http.Handler) bool {
	return srv.AddHttpRouterHandler(method, true, url, NewHttpRouterHandle(handler))
}

func (srv *BaseServer) AddServerHandlerFunc(method string, url string, handlerFunc func(http.ResponseWriter, *http.Request)) bool {
	return srv.AddHttpRouterHandler(method, true, url, NewHttpRouterHandle(http.HandlerFunc(handlerFunc)))
}



