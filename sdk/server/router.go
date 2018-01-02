package server

import (
	"path/filepath"
	"github.com/julienschmidt/httprouter"
)

//new router set httprouter options and static dir
func newRouter(assetLoader IAssetLoader) *httprouter.Router{
	router := httprouter.New()

	//disable automatic redirection if the current route can't be matched but a
	// handler for the path with (without) the trailing slash exists.
	router.RedirectTrailingSlash = false

	assetPath := filepath.Join("/", assetLoader.GetRootDirname(), "*filepath")
	//router.ServeFiles(assetPath, http.FileServer(assetLoader.GetHttpFilesystem()))
	router.ServeFiles(assetPath, assetLoader.GetHttpFilesystem())

	return router
}
