package server

import (
	"path/filepath"
)

type BaseAssetLoader struct {
	RootDir string
	RootUrl	string
}

//GetAssetPath: get asset path based on asset name, it joins the asset root with relative asset path
func (loader *BaseAssetLoader) GetAssetPath(assetName string) string {
	return filepath.Join(loader.RootDir, loader.getRelAssetPath(assetName))
}

func (loader *BaseAssetLoader) GetRootUrl() string {
	return loader.RootUrl
}

//getRelAssetPath: get relative asset path based on asset name
func (loader *BaseAssetLoader) getRelAssetPath(assetName string) string {
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