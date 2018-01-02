package server

import (
	"path/filepath"
	"net/http"
	"log"
)

type IAssetLoader interface{
	//get loader's name
	GetName()									string
	//get asset root dirname
	GetRootDirname()							string
	//asset root url which will be used in template
	GetRootUrl(namespace string)				string
	//used for http static dir
	GetHttpFilesystem()        					http.FileSystem
	//check asset path exists or not
	Exists(assetPath string) 					bool
	// get asset path
	GetAssetPath(namespace, assetName string)	string
	// read asset data
	ReadBytes(assetPath string)  				([]byte, error)
}

type BaseLoader struct {
	PkgRoot 		string
	PkgRootUrl		string
	VendorRoot		string
	VendorRootUrl   string
}

//newAssetLoader: there are two asset loader type now:
//1. go-bindata: generate data object file by go-bindata, and read asset data from go object file
//2. localfs: read asset data from local filesystem
//if using go-bindata, it need pass in go-bindata funcs: Asset,AssetDir,AssetInfo
//which will be used to generate the pseudo http filesystem to serve static files
func newAssetLoader(serverName, assetDir string, args... interface{}) IAssetLoader {
	if assetDir != "" {
		return NewLocalFsLoader(serverName, assetDir)
	}

	//if not specify assetDir, then we assume it is using go-bindata internal asset
	if len(args) != 3 {
		log.Println("Insufficient bindata funcs")
		return nil
	}

	//try get bindata related functions from args
	assetFunc, assetDirFunc, assetInfoFunc := getBindataFuncs(args...)
	if assetFunc == nil || assetDirFunc == nil || assetInfoFunc == nil {
		log.Println("bindata asset related funcs cannot be nil")
		return nil
	}

	return NewBindataLoader(serverName, assetFunc, assetDirFunc, assetInfoFunc)
}

//GetRootUrl: get asset root url which will be used in template files based on namespace
//now only support two kinds of asset namespace:
//1. pkg:    <ASSET_ROOT>/pkg/<PKG_NAME>
//2. vendor: <ASSET_ROOT/vendor
//
//     |- pkg ---------------------------------
//         |- deploy -------------------------| /asset/pkg/deploy: ASSET_ROOT_URL->/ASSET_ROOTDIR_NAME/NAMESPACE/PKG_NAME
//               |- js------------------------|
//                   * main.js ---------------| js/main.js: REL_ASSET_DIR
//    |- vendor
//         |- bootstrap
//                |- js
//                    |- boostrap.min.js
func (loader *BaseLoader) GetRootUrl(namespace string) string {
	switch namespace {
	case "pkg":
		return loader.PkgRootUrl
	case "vendor":
		return loader.VendorRootUrl
	default:
		log.Println("Non supported namespace")
		return ""
	}
}

//GetAssetPath: get asset path based on  asset name, it joins the asset root with relative asset path
func (loader *BaseLoader) GetAssetPath(namespace, assetName string) string {
	var assetPath string

	switch namespace {
	case "pkg":
		assetPath = filepath.Join(loader.PkgRoot, loader.getRelAssetPath(assetName))
	case "vendor":
		assetPath = filepath.Join(loader.VendorRoot, loader.getRelAssetPath(assetName))
	default:
		log.Println("Non supported namespace")
	}

	return assetPath
}

//getRelAssetPath: get relative asset path based on asset name
func (loader *BaseLoader) getRelAssetPath(assetName string) string {
	extension := filepath.Ext(assetName)
	switch extension {
	case ".html",".htm":
		return filepath.Join("html", assetName)
	case ".js":
		return filepath.Join("js", assetName)
	case ".css":
		return filepath.Join("css", assetName)
	case ".yaml",".yml":
		return filepath.Join("yaml", assetName)
	default:
		return filepath.Join("misc", assetName)
	}
}

//initRootInfo: for performance consideration,
//we generate the dedicate root for different kinds of assets
//when each asset load initializing
func (loader *BaseLoader) initRootInfo(rootDirname, serverName string) {
	loader.PkgRoot = filepath.Join(rootDirname, "pkg", serverName)
	loader.PkgRootUrl = filepath.Join("/", loader.PkgRoot)

	loader.VendorRoot = filepath.Join(rootDirname, "vendor")
	loader.VendorRootUrl = filepath.Join("/", loader.VendorRoot)
}
