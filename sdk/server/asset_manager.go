package server

import (
	"log"
	"net/http"
)

type IAssetManager interface {
	GetLoaderType()							string
	GetRootUrl(namespace string)			string
	GetHttpFilesystem(namespace string)		http.FileSystem

	GetAssetPath(assetName string)          string
	ReadBytes(assetPath string) 			([]byte, error)
}

//try get assetdir from args, normally it should be the first argument
func getAssetDirFromArgs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}

	firstArg := args[0]
	assetDir, ok := firstArg.(string)
	if !ok {
		return ""
	}

	return assetDir
}

//newAssetManager: there are two asset loader type now:
//1. go-bindata: generate data object file by go-bindata, and read asset data from go object file
//2. localfs: read asset data from local filesystem
//if using go-bindata, it need pass in go-bindata funcs: Asset,AssetDir,AssetInfo
//which will be used to generate the pseudo http filesystem to serve static files
func newAssetManager(serverName string, args ...interface{}) IAssetManager {
	assetDir := getAssetDirFromArgs(args...)
	//if we specified assetDir, then all loader should load asset from such dir
	if assetDir != "" {
		return NewLocalFsAssetManager(serverName, assetDir)
	}

	//if not specify assetDir, then we assume it is using go-bindata internal asset
	//There should be 4 items in slice: assetDir, Asset, AssetDir, AssetInfo
	if len(args) != 4 {
		log.Println("Insufficient bindata funcs")
		return nil
	}

	//try get package specific bindata related functions from args
	remainArgs := args[1:]
	pkgAssetFunc, pkgAssetDirFunc, pkgAssetInfoFunc := getBindataFuncs(remainArgs...)
	if pkgAssetFunc == nil || pkgAssetDirFunc == nil || pkgAssetInfoFunc == nil {
		log.Println("bindata asset related funcs cannot be nil")
		return nil
	}

	return NewBindataManager(serverName, pkgAssetFunc, pkgAssetDirFunc, pkgAssetInfoFunc)
}






