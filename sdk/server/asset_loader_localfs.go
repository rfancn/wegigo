package server

import (
	"net/http"
	"path/filepath"
	"io/ioutil"
	"os"
)

//External Asset generate By local filesystem
type LocalFsLoader struct {
	BaseLoader
	//the assetDir option specified when running wegigo
	AssetDir string
}

func NewLocalFsLoader(serverName, assetDir string) *LocalFsLoader {
	loader := &LocalFsLoader{AssetDir: assetDir}
	loader.initRootInfo(loader.GetRootDirname(), serverName)

	return loader
}

func (loader *LocalFsLoader) GetName() string {
	return "localfs"
}

func (loader *LocalFsLoader) GetHttpFilesystem() http.FileSystem {
	return http.Dir(loader.AssetDir)
}

//root dirname should be the specified assetDir's last dirname
func (loader *LocalFsLoader) GetRootDirname() string {
	return filepath.Base(loader.AssetDir)
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

