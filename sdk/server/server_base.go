package server

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"github.com/kabukky/httpscerts"
	"github.com/julienschmidt/httprouter"
	"path/filepath"
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

func (srv *BaseServer) AddRoute(method string, url string, handle httprouter.Handle) bool {
	lowerMethod := strings.ToLower(method)
	absUrl := filepath.Join("/", srv.Name, url)
	switch lowerMethod {
	case "get":
		srv.router.GET(absUrl, handle)
	case "post":
		srv.router.POST(absUrl, handle)
	default:
		log.Println("Non supported http method: ", method)
		return false
	}

	return true
}

