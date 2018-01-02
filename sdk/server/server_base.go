package server

import (
	"fmt"
	"net/http"
	"strings"
	"log"
	"github.com/kabukky/httpscerts"
	"github.com/julienschmidt/httprouter"
	"github.com/flosch/pongo2"
)

type BaseServer struct {
	Name        string
	router      *httprouter.Router
	assetLoader IAssetLoader
	templateSet *pongo2.TemplateSet
}

func (srv *BaseServer) Initialize(serverName, assetDir string, args... interface{}) bool {
	srv.Name = serverName

	srv.assetLoader = newAssetLoader(serverName, assetDir, args...)
	if srv.assetLoader == nil {
		log.Println("Failed to initialize asset loader")
		return false
	}

	srv.router = newRouter(srv.assetLoader)
	if srv.router == nil {
		log.Println("Failed to initialize server router")
		return false
	}

	srv.templateSet = newTemplateSet(srv.assetLoader)
	if srv.templateSet == nil {
		log.Println("Failed to initialize server template set")
		return false
	}

	return true
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
	switch lowerMethod {
	case "get":
		srv.router.GET(url, handle)
	case "post":
		srv.router.POST(url, handle)
	default:
		log.Println("Non supported http method: ", method)
		return false
	}

	return true
}

func (srv *BaseServer) GetAssetPath(assetName string, args ...interface{}) string {
	//by default, namespace is "pkg" unless you specify the namespace in args
	namespace := ASSET_NAMESPACE_PKG
	if len(args) == 1 {
		namespace = args[0].(string)
	}

	return srv.assetLoader.GetAssetPath(namespace, assetName)
}


func (srv *BaseServer) ReadAssetBytes(assetName string, args ...interface{})	([]byte, error) {
	assetPath := srv.GetAssetPath(assetName, args...)
	return srv.assetLoader.ReadBytes(assetPath)
}

func (srv *BaseServer) GetAssetLoaderName()	 string {
	return srv.assetLoader.GetName()
}


