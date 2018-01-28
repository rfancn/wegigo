package server

import (
	"path/filepath"
)

type BindataAssetManager struct {
	BaseAssetManager
}

func NewBindataManager(pkgName string, pkgAssetFunc BindataFuncAsset, pkgAssetDirFunc BindataFuncAssetDir, pkgAssetInfoFunc BindataFuncAssetInfo) *BindataAssetManager {
	manager := &BindataAssetManager{}

	manager.loaderType = ASSET_LOADER_TYPE_BINDATA

	//new pkg asset vendor asset loader
	pkgAssetRelPath := filepath.Join(ASSET_NAMESPACE_PKG, pkgName)
	manager.pkgAssetLoader = NewBindataLoader(pkgAssetRelPath, pkgAssetFunc, pkgAssetDirFunc, pkgAssetInfoFunc)
	//new bindata vendor asset loader
	manager.vendorAssetLoader = NewBindataLoader(ASSET_NAMESPACE_VENDOR, Asset, AssetDir, AssetInfo)

	return manager
}
