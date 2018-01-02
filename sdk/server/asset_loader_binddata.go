package server

import (
	"os"
	"net/http"
	"reflect"
	"github.com/elazarl/go-bindata-assetfs"
)

var BINDATA_ASSET_ROOT_DIRNAME = "asset"

//abbreviation of the important go-bindata func
type BindataFuncAsset func(string) ([]byte, error)
type BindataFuncAssetDir  func(string) ([]string, error)
type BindataFuncAssetInfo func(string) (os.FileInfo, error)

type BindataLoader struct {
	BaseLoader
	AssetFunc  		BindataFuncAsset
	AssetDirFunc   	BindataFuncAssetDir
	AssetInfoFunc  	BindataFuncAssetInfo
}

//try locate the go-bindata funcs by passed in arguments
//in this way, the order of the bindata funcs passed in are not sensitive
func getBindataFuncs(args... interface{}) (assetFunc BindataFuncAsset, assetDirFunc BindataFuncAssetDir, assetInfoFunc BindataFuncAssetInfo){
	for _, arg := range args {
		t := reflect.TypeOf(arg)
		if t.Kind() == reflect.Func {
			tName := t.String()
			switch tName {
			//BindataFuncAsset, as byte is uint8 type, need use uint8 here
			case "func(string) ([]uint8, error)":
				assetFunc = arg.(func(string) ([]uint8, error))
			//BindataFuncAssetDir
			case "func(string) ([]string, error)":
				assetDirFunc = arg.(func(string) ([]string, error))
			//BindataFuncAssetInfo
			case "func(string) (os.FileInfo, error)":
				assetInfoFunc = arg.(func(string) (os.FileInfo, error))
			}
		}
	}

	return
}

//NewBindataLoader: the last 3 one are functions exported by go-bindata
func NewBindataLoader(serverName string, assetFunc	BindataFuncAsset, assetDirFunc BindataFuncAssetDir, assetInfoFunc BindataFuncAssetInfo) *BindataLoader {
	loader := &BindataLoader{AssetFunc: assetFunc, AssetDirFunc: assetDirFunc, AssetInfoFunc: assetInfoFunc}
	//for performance consideration,
	//it need save the root info for future reference
	loader.initRootInfo(loader.GetRootDirname(), serverName)

	return loader
}

func (loader *BindataLoader) GetName() string {
	return "bindata"
}

//In project source files, go-bindata always use "asset" as the root dirname
func (loader *BindataLoader) GetRootDirname() string {
	return BINDATA_ASSET_ROOT_DIRNAME
}

//Generate http filesystem for sever http static files
func (loader *BindataLoader) GetHttpFilesystem() http.FileSystem {
	//prefix must be the asset Rootdir name,
	//assetfs library will combine the assetRootDirname + relpath of the asset name
	prefix := loader.GetRootDirname()
	return &assetfs.AssetFS{Prefix: prefix, Asset: loader.AssetFunc, AssetDir: loader.AssetDirFunc,AssetInfo: loader.AssetInfoFunc}
}

//Check specified asset path exist or not
func (loader *BindataLoader) Exists(assetPath string) bool {
	if _, err := loader.AssetInfoFunc(assetPath); err == nil {
		return true
	}

	return false
}

//Read asset data
func (loader *BindataLoader) ReadBytes(assetPath string) ([]byte, error) {
	data, err := loader.AssetFunc(assetPath)
	if err != nil {
		return nil, err
	}

	return data, nil
}


