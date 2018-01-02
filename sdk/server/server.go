package server

import (
	"net/http"
)

type IServer interface {
	Initialize(name, assetDir string)					bool
	//ReadAssetBytes read the asset data from "pkg" namespace by default,
	//if you want to read it from "vendor" namespace, then specify namespace in args
	ReadAssetBytes(assetName string, args... interface{})		([]byte, error)
	RespRender(w http.ResponseWriter, templatePath string, context map[string]interface{})
}


