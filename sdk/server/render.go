package server

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"github.com/flosch/pongo2"
)

//new pongo2 template set based on loader instance's type
//there are two template set exist at this moment
//1. memory loader:  pass bindataLoader as arg, and read template data from golang object
//2. localfs loader: used pongo2 internal localfs loader
func newTemplateSet(assetManager IAssetManager) (tSet *pongo2.TemplateSet) {
	loaderType := assetManager.GetLoaderType()
	switch loaderType {
	case ASSET_LOADER_TYPE_BINDATA:
		memLoader := NewMemoryTemplateLoader(assetManager)
		tSet = pongo2.NewSet("memory", memLoader)
	case ASSET_LOADER_TYPE_LOCALFS:
		fileLoader := pongo2.MustNewLocalFileSystemLoader("")
		tSet = pongo2.NewSet("localfs", fileLoader)
	default:
		log.Println("Invalid server asset type")
		tSet = nil
	}

	return tSet
}

//Response and render templates
func (srv *SimpleServer) Render(templatePath string, context map[string]interface{}) (string, error) {
	if context == nil {
		context = make(map[string]interface{})
	}

	//add root url for template reference
	context["PKG_ROOT"] = srv.assetManager.GetRootUrl(ASSET_NAMESPACE_PKG)
	context["VENDOR_ROOT"] = srv.assetManager.GetRootUrl(ASSET_NAMESPACE_VENDOR)

	t := pongo2.Must(srv.templateSet.FromFile(templatePath))
	output, err := t.Execute(context)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return output, nil
}

//Response and render generic
func (srv *SimpleServer) RespRender(w http.ResponseWriter, templateFile string, context map[string]interface{}) bool {
	templatePath := srv.assetManager.GetAssetPath(templateFile)
	if templatePath == "" {
		log.Println("Error locate template:", templateFile)
		return false
	}

	output, err := srv.Render(templatePath, context)
	if err != nil {
		fmt.Printf("Error render template[%s]: %s\n", templatePath, err)
		return false
	}

	_, err = fmt.Fprint(w, output)
	if err != nil {
		fmt.Printf("Error echo server response: ", err)
		return false
	}

	return true
}

func (srv *SimpleServer) RespText(w http.ResponseWriter, response string) bool {
	_, err := fmt.Fprint(w, response)
	if err != nil {
		fmt.Printf("Error return text response: ", err)
		return false
	}

	return true
}

func (srv *SimpleServer) RespJson(w http.ResponseWriter, response interface{}) bool {
	w.Header().Set("Content-Type", "application/json")

	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}

	_, err = fmt.Fprint(w, string(js))
	if err != nil {
		fmt.Printf("Error return json response: ", err)
		return false
	}

	return true
}

