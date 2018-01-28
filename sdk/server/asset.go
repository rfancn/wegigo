package server

import (
	"path/filepath"
	"os"
	"log"
)

const ASSET_NAMESPACE_PKG = "pkg"
const ASSET_NAMESPACE_VENDOR = "vendors"

var ASSET_NAMESPACES = []string {
	ASSET_NAMESPACE_PKG,
	ASSET_NAMESPACE_VENDOR,
}

const ASSET_LOADER_TYPE_LOCALFS = "localfs"
const ASSET_LOADER_TYPE_BINDATA = "go-bindata"


//getBasedir returns absolute wegigo running dir
func getBasedir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("Error get current execute path: %s", err)
		return ""
	}

	return dir
}


