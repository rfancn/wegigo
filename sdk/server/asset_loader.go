package server

import "net/http"

type IAssetLoader interface{
	//get asset root dirname
	GetRootDirname()						string
	GetRootUrl() 							string
	//used for http static dir
	GetHttpFilesystem()        				http.FileSystem
	//check asset path exists or not
	Exists(assetPath string) 				bool
	// read asset data
	ReadBytes(assetPath string)  			([]byte, error)

	//get asset path, implemented in BaseAssetLoader
	GetAssetPath(assetName string)			string
}

