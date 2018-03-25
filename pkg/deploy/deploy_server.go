package deploy

import (
	"log"
	"github.com/rfancn/wegigo/sdk/server"
	"github.com/rfancn/goy2h"
	"net/http"
	"io/ioutil"
	"os"
	"io"
)

const SERVER_NAME = "deploy"

type ServerStatus int

const (
	STATUS_INIT ServerStatus = iota
	STATUS_CONFIG_STARTED
	STATUS_CONFIG_DONE
	STATUS_DEPLOY_STARTED
)

type CmdArgument struct {
	ServerUrl string
	AssetDir  string
	Timeout   int
}

type DeployServer struct {
	server.SimpleServer
	cmdArg      *CmdArgument
	y2h 		*goy2h.Y2H
	status  	ServerStatus
	response    *ServerJsonResponse
	timeout 	int
}

type ServerJsonResponse struct {
	Result string
	Detail string
}

func NewDeployServer(serverName string, cmdArg *CmdArgument) *DeployServer {
	srv := &DeployServer{}

	//when initialize simpleserver, the first argument must be assetDir
	if ! srv.SimpleServer.Initialize(serverName, cmdArg.AssetDir, Asset, AssetDir, AssetInfo) {
		return nil
	}

	//set status to be init at the very beginning
	srv.status = STATUS_INIT
	srv.response = &ServerJsonResponse{Result:"", Detail:""}

	srv.y2h = goy2h.New()

	return srv
}

func Run(cmdArg *CmdArgument) {
	log.Println("Run deploy Server at:", cmdArg.ServerUrl)

	srv := NewDeployServer(SERVER_NAME, cmdArg)
	if srv == nil {
		log.Fatal("Error create deploy server")
	}

	//get deploy related files
	srv.prepareDeployEnv()

	srv.setupRouter()

	err := srv.Run(cmdArg.ServerUrl)
	if err != nil {
		log.Fatal("Error start deploy server:", err)
	}

}

func (srv *DeployServer) setupRouter() {
	srv.AddServerHandlerFunc("get", "/", srv.ViewIndex)
	srv.AddServerHandlerFunc("post", "/config", srv.ViewConfig)
	srv.AddServerHandlerFunc("get", "/deploy", srv.ViewDeploy)
	srv.AddServerHandlerFunc("get", "/run", srv.ViewRun)
}

//restoreFromBindata: read bindata object content and save to target file
func (srv *DeployServer) restoreFromBindata(srcFile, targetFile string) bool {
	content, err := srv.ReadAssetBytes(srcFile)
	if err != nil {
		log.Println("Error get from bindata objects:", srcFile)
		return false
	}
	err = ioutil.WriteFile(targetFile, content, 0664)
	if err != nil {
		log.Println("Error prepare deploy file:", targetFile)
		return false
	}

	log.Println("Generated deploy file from bindata object:", targetFile)
	return true
}

//copyFromLocalfs: copy src file to target file
func (srv *DeployServer) copyFromLocalfs(srcFile, targetFile string) bool {
	srcFilepath := srv.GetAssetPath(srcFile)
	in, err := os.Open(srcFilepath)
	if err != nil {
		log.Println("copyFromLocalfs(): Error open file:", srcFilepath)
		return false
	}
	defer in.Close()

	out, err := os.Create(targetFile)
	if err != nil {
		log.Println("copyFromLocalfs(): Error create target file:", targetFile)
		return false
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		log.Println("copyFromLocalfs(): Error copy file content:", err)
		return false
	}
	return true
}


func (srv *DeployServer) prepareDeployEnv() {
	assetLoaderType := srv.GetLoaderType()
	deployFiles := map[string]string{"deploy.sh":"deploy.sh", "kube_setup.playbook":"kube_setup.yaml"}

	for srcFile, targetFile:= range deployFiles {
		switch assetLoaderType {
		case server.ASSET_LOADER_TYPE_BINDATA:
			srv.restoreFromBindata(srcFile, targetFile)
		case server.ASSET_LOADER_TYPE_LOCALFS:
			srv.copyFromLocalfs(srcFile, targetFile)
		}
	}
}

func (srv *DeployServer) RespJsonSuccess(w http.ResponseWriter, detail string) {
	srv.response.Result = "success"
	srv.response.Detail = detail
	srv.RespJson(w, srv.response)
}

func (srv *DeployServer) RespJsonError(w http.ResponseWriter, detail string) {
	srv.response.Result = "error"
	srv.response.Detail = detail
	srv.RespJson(w, srv.response)
}

