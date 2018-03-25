package deploy

import (
	"net/http"
	"log"
	"github.com/gorilla/websocket"
)

func internalError(ws *websocket.Conn, msg string, err error) {
	log.Println(msg, err)
	ws.WriteMessage(websocket.TextMessage, []byte("Internal server error."))
}

func (srv *DeployServer) ViewDeploy(w http.ResponseWriter, r *http.Request) {
	log.Println("Enter ViewDeploy, server status is:", srv.status)

	/**
	if srv.status != STATUS_CONFIG_DONE {
		log.Println("Server seems not get configured!")
		http.Error(w,"You must configure deployment options firstly!", http.StatusInternalServerError)
		return
	}

	if srv.status == STATUS_DEPLOY_STARTED {
		log.Println("Receive duplicate deploy request at the same time.")
		http.Error(w,"There is a deploy routine already running!", http.StatusInternalServerError)
		return
	}
	**/

	srv.RespRenderFile(w, "deploy.html", nil)

	log.Println("Exit ViewDeploy, server status is:", srv.status)
}


func (srv *DeployServer) ViewRun(w http.ResponseWriter, r *http.Request) {
	log.Println("Enter ViewRun, server status is:", srv.status)

	ws := newWebsocket(w, r)
	defer ws.Close()
	if ws == nil {
		http.Error(w, "Error new websocket server", http.StatusInternalServerError)
		return
	}

	runServerCommand(ws, "/bin/sh", "deploy.sh")

	//closeWebsocket(ws)

	log.Println("Exit ViewRun, server status is:", srv.status)
}
