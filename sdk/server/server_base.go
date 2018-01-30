package server

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"github.com/kabukky/httpscerts"
	"github.com/julienschmidt/httprouter"
	"path/filepath"
	"context"
)

type BaseServer struct {
	Name        string
	router      *httprouter.Router
}

func (srv *BaseServer) Initialize(serverName string, args ...interface{}) bool {
	srv.Name = serverName
	srv.router = httprouter.New()

	return true
}

func (srv *BaseServer) ResetHttpRouter() {
	srv.router = httprouter.New()
}

func (srv *BaseServer) RunHttp(bind string, port int) error {
	listen := fmt.Sprintf("%s:%d", bind, port)
	return http.ListenAndServe(listen, srv.router)
}

func (srv *BaseServer) RunHttps(bind string, port int) error {
	// Check if the cert files are available.
	if err := httpscerts.Check("cert.pem", "key.pem"); err != nil {
		// If they are not available, generate new ones.
		if err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:443"); err != nil {
			log.Println("Failed to generate secure credentials")
			return err
		}
	}

	listen := fmt.Sprintf("%s:%d", bind, port)
	return http.ListenAndServeTLS(listen, "cert.pem", "key.pem", srv.router)
}

//add httprouter handle func
func (srv *BaseServer) AddHttpRouterHandler(method string, url string, handleFunc httprouter.Handle) bool {
	lowerMethod := strings.ToLower(method)
	absUrl := filepath.Join("/", srv.Name, url)
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

//add http handle func
func (srv *BaseServer) AddHttpHandler(method string, url string, handler http.Handler) bool {
	return srv.AddHttpRouterHandler(method, url, NewHttpRouterHandle(handler))
}

func (srv *BaseServer) AddHttpHandlerFunc(method string, url string, handlerFunc func(http.ResponseWriter, *http.Request)) bool {
	return srv.AddHttpRouterHandler(method, url, NewHttpRouterHandle(http.HandlerFunc(handlerFunc)))
}


