package server

import (
	"fmt"
	"net/http"
	"github.com/flosch/pongo2"
	"log"
	"encoding/json"
)

//new pongo2 template set based on loader instance's type
//there are two template set exist at this moment
//1. memory loader:  pass bindataLoader as arg, and read template data from golang object
//2. localfs loader: used pongo2 internal localfs loader
func newTemplateSet(assetLoader IAssetLoader) (tSet *pongo2.TemplateSet) {
	loaderName := assetLoader.GetName()
	switch loaderName {
	case "bindata":
		memLoader := NewMemoryTemplateLoader(assetLoader)
		tSet = pongo2.NewSet("memory", memLoader)
	case "localfs":
		fileLoader := pongo2.MustNewLocalFileSystemLoader("")
		tSet = pongo2.NewSet("localfs", fileLoader)
	default:
		log.Println("Invalid server asset type")
		tSet = nil
	}

	return tSet
}

func newTemplateSet123(assetLoader IAssetLoader) (tSet *pongo2.TemplateSet) {
	loaderType := getInterfaceType(assetLoader)
	switch loaderType {
	case "BindataLoader":
		memLoader := NewMemoryTemplateLoader(assetLoader)
		tSet = pongo2.NewSet("memory", memLoader)
	case "LocalFsLoader":
		fileLoader := pongo2.MustNewLocalFileSystemLoader("")
		tSet = pongo2.NewSet("localfs", fileLoader)
	default:
		log.Println("Invalid server asset type")
		tSet = nil
	}

	return tSet
}

//template lookup path will be checked by following order:
//- from pkg asset dir
//- from vendor asset dir
func (srv *BaseServer) LookupTemplate(templateFile string) string {
	var namespaces = []string{"pkg", "vendor"}

	for _, nm := range namespaces {
		templatePath := srv.assetLoader.GetAssetPath(nm, templateFile)
		if srv.assetLoader.Exists(templatePath) {
			return templatePath
		}
	}

	return ""
}

//Response and render templates
func (srv *BaseServer) Render(templatePath string, context map[string]interface{}) (string, error) {
	if context == nil {
		context = make(map[string]interface{})
	}

	//add root url for template reference
	context["PKG_ROOT"] = srv.assetLoader.GetRootUrl("pkg")
	context["VENDOR_ROOT"] = srv.assetLoader.GetRootUrl("vendor")

	t := pongo2.Must(srv.templateSet.FromFile(templatePath))
	output, err := t.Execute(context)
	if err != nil {
		log.Println(err)
		return "", err
	}

	return output, nil
}

//Response and render generic
func (srv *BaseServer) RespRender(w http.ResponseWriter, templateFile string, context map[string]interface{}) bool {
	templatePath := srv.LookupTemplate(templateFile)
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

func (srv *BaseServer) RespText(w http.ResponseWriter, response string) bool {
	_, err := fmt.Fprint(w, response)
	if err != nil {
		fmt.Printf("Error return text response: ", err)
		return false
	}

	return true
}

func (srv *BaseServer) RespJson(w http.ResponseWriter, response interface{}) bool {
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

