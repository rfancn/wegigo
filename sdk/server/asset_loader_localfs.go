package server

import (
	"net/http"
	"path/filepath"
	"io/ioutil"
	"os"
)

//External Asset generate By local filesystem
type LocalFsLoader struct {
	BaseAssetLoader
	//the assetDir option specified when running wegigo
	AssetDir string
}

// assetDir is the whole asset root dir,
// divide asset into two logic types of asset:
// 1. vendor asset is vendor specific asset, like 3rd party js/css libraries
// 2. pkg asset belongs to specific package, according to codes in src/pkg/<pkg name>, it is normally serverName
func NewLocalFsLoader(assetDir string, relPath string) *LocalFsLoader {
	loader := &LocalFsLoader{AssetDir: assetDir}

	rootDirname := loader.GetRootDirname()
	loader.RootDir = filepath.Join(rootDirname, relPath)
	loader.RootUrl = filepath.Join("/", loader.RootDir)

	return loader
}

//root dirname should be the specified assetDir's last dirname
func (loader *LocalFsLoader) GetRootDirname() string {
	return filepath.Base(loader.AssetDir)
}

func (loader *LocalFsLoader) GetHttpFilesystem() http.FileSystem {
	return http.Dir(loader.RootDir)
}

func (loader *LocalFsLoader) Exists(assetPath string) bool {
	if _, err := os.Stat(assetPath); err == nil {
		return true
	}

	return false
}

func (loader *LocalFsLoader) ReadBytes(assetPath string) ([]byte, error) {
	content, err := ioutil.ReadFile(assetPath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

