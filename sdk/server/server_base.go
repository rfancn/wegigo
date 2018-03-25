package server

import (
	"net/http"
	"net/url"
	"log"
	"github.com/kabukky/httpscerts"
	"github.com/julienschmidt/httprouter"

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

//Run Server based on serverUrl
func (srv *BaseServer) Run(serverUrl string) error {
	u, err := url.Parse(serverUrl)
	if err != nil {
		return err
	}

	//u.Host is a host:port string
	switch u.Scheme {
	case "http":
		return srv.RunHttp(u.Host)
	case "https":
		return srv.RunHttps(u.Host)
	}

	return nil
}

//Run http server
func (srv *BaseServer) RunHttp(listen string) error {
	return http.ListenAndServe(listen, srv.router)
}

//Run https server
func (srv *BaseServer) RunHttps(listen string) error {
	// Check if the cert files are available.
	if err := httpscerts.Check("cert.pem", "key.pem"); err != nil {
		// If they are not available, generate new ones.
		if err = httpscerts.Generate("cert.pem", "key.pem", "127.0.0.1:443"); err != nil {
			log.Println("Failed to generate secure credentials")
			return err
		}
	}

	return http.ListenAndServeTLS(listen, "cert.pem", "key.pem", srv.router)
}
