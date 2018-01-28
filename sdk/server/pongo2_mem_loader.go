//copy from: https://github.com/sharpner/pobin/
package server

import (
	"bytes"
	"io"
	"path/filepath"
	"github.com/flosch/pongo2"
)

type LoadFunc func(path string) ([]byte, error)

type MemoryTemplateLoader struct {
	assetManager IAssetManager
}

//NewMemoryTemplateLoader loads a go-bindata object data
func NewMemoryTemplateLoader(assetManager IAssetManager) pongo2.TemplateLoader {
	return &MemoryTemplateLoader{assetManager: assetManager}
}

// Abs resolves a filename relative to the base directory. Absolute paths are allowed.
// When there's no base dir set, the absolute path to the filename
// will be calculated based on either the provided base directory (which
// might be a path of a template which includes another template) or
// the current working directory.
func (m MemoryTemplateLoader) Abs(base, name string) string {
	if filepath.IsAbs(name) || base == "" {
		return name
	}

	if name == "" {
		return base
	}

	return filepath.Dir(base) + string(filepath.Separator) + name
}

// Get reads the path's content from your local filesystem.
func (m MemoryTemplateLoader) Get(path string) (io.Reader, error) {
	data, err := m.assetManager.ReadBytes(path)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
