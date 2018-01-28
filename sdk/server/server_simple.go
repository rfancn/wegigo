package server

import (
	"path/filepath"
	"log"
	"github.com/flosch/pongo2"
)

type SimpleServer struct {
	BaseServer
	assetManager IAssetManager
	templateSet *pongo2.TemplateSet
}


//when initialize simpleserver, the first argument must be assetDir
//normally, it should be: assetDir, Asset, AssetDir, AssetInfo
func (srv *SimpleServer) Initialize(serverName string, args ...interface{}) bool {
	srv.BaseServer.Initialize(serverName, args...)

	srv.assetManager = newAssetManager(serverName, args...)
	if srv.assetManager == nil {
		log.Println("Failed to initialize asset manager")
		return false
	}

	srv.setupRouter(srv.assetManager)

	srv.templateSet = newTemplateSet(srv.assetManager)
	if srv.templateSet == nil {
		log.Println("Failed to initialize server template set")
		return false
	}

	return true
}

//new router set httprouter options and static dir
func (srv *SimpleServer) setupRouter(assetManager IAssetManager) {
	//disable automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	//srv.router.RedirectTrailingSlash = false

	for _, nm := range ASSET_NAMESPACES {
		url := filepath.Join(assetManager.GetRootUrl(nm), "*filepath")
		srv.router.ServeFiles(url, assetManager.GetHttpFilesystem(nm))
	}
}

/**
 ** Asset Related API
 */
func (srv *SimpleServer) ReadAssetBytes(assetName string)	([]byte, error) {
	assetPath := srv.GetAssetPath(assetName)
	return srv.assetManager.ReadBytes(assetPath)
}

func (srv *SimpleServer) GetAssetPath(assetName string) string {
	return srv.assetManager.GetAssetPath(assetName)
}

func (srv *SimpleServer) GetLoaderType() string {
	return srv.assetManager.GetLoaderType()
}



