package server

import (
	"path/filepath"
	"os"
	"log"
)


var ASSET_NAMESPACE_PKG = "pkg"
var ASSET_NAMESPACE_VENDOR = "vendor"

//getBasedir returns absolute wegigo running dir
func getBasedir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("Error get current execute path: %s", err)
		return ""
	}

	return dir
}


