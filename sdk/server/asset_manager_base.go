package server

import (
	"log"
	"net/http"
)

type BaseAssetManager struct {
	loaderType  string
	pkgAssetLoader IAssetLoader
	vendorAssetLoader IAssetLoader
}

func (m *BaseAssetManager) GetAssetPath(assetName string) string {
	assetPath := m.pkgAssetLoader.GetAssetPath(assetName)
	if assetPath == "" {
		assetPath = m.vendorAssetLoader.GetAssetPath(assetName)
	}

	return assetPath
}

func (m *BaseAssetManager) GetLoaderType() string {
	return m.loaderType
}

//GetRootUrl: get asset root url which will be used in template files based on namespace
//now only support two kinds of asset namespace:
//1. pkg:    <ASSET_ROOT>/pkg/<PKG_NAME>
//2. vendor: <ASSET_ROOT/vendor
//
//     |- pkg ---------------------------------
//         |- deploy -------------------------| /asset/pkg/deploy: ASSET_ROOT_URL->/ASSET_ROOTDIR_NAME/NAMESPACE/PKG_NAME
//               |- js------------------------|
//                   * main.js ---------------| js/main.js: REL_ASSET_DIR
//    |- vendor
//         |- bootstrap
//                |- js
//                    |- boostrap.min.js
func (m *BaseAssetManager) GetRootUrl(namespace string) string {
	switch namespace {
	case ASSET_NAMESPACE_PKG:
		return m.pkgAssetLoader.GetRootUrl()
	case ASSET_NAMESPACE_VENDOR:
		return m.vendorAssetLoader.GetRootUrl()
	default:
		log.Println("Non supported namespace")
		return ""
	}
}

func (m *BaseAssetManager) GetHttpFilesystem(namespace string) http.FileSystem {
	switch namespace {
	case ASSET_NAMESPACE_PKG:
		return m.pkgAssetLoader.GetHttpFilesystem()
	case ASSET_NAMESPACE_VENDOR:
		return m.vendorAssetLoader.GetHttpFilesystem()
	default:
		log.Println("AssetManager GetHttpFilesystem(): non-supported namespace")
		return nil
	}
}

func (m *BaseAssetManager) ReadBytes(assetPath string) ([]byte,error) {
	//try read from package asset, then try from vendor asset
	data, err := m.pkgAssetLoader.ReadBytes(assetPath)
	if err != nil {
		data, err = m.vendorAssetLoader.ReadBytes(assetPath)
	}

	return data, err
}
