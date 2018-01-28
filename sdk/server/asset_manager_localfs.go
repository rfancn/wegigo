package server

import (
	"path/filepath"
)

type LocalFsAssetManager struct {
	BaseAssetManager
}

func NewLocalFsAssetManager(pkgName, assetDir string) IAssetManager{
	manager := &LocalFsAssetManager{}

	manager.loaderType = ASSET_LOADER_TYPE_LOCALFS

	//new pkg asset loader
	pkgAssetRelPath := filepath.Join(ASSET_NAMESPACE_PKG, pkgName)
	manager.pkgAssetLoader = NewLocalFsLoader(assetDir, pkgAssetRelPath)

	//new vendor asset loader, relpath is ASSET_NAMESPACE_VENDOR
	manager.vendorAssetLoader = NewLocalFsLoader(assetDir, ASSET_NAMESPACE_VENDOR)

	return manager
}
