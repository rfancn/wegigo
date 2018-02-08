package deploy

import (
	"net/http"
	"bytes"
	"log"
	"github.com/julienschmidt/httprouter"
)

var INSTALL_WIZARD = map[string]string {
	"general": 	"step_general.yaml",
	"wechat":     	"step_wechat.yaml",
	"server":   	"step_server.yaml",
	"summary":   	"step_summary.yaml",
}

//yaml definition for edit server
const DEPLOY_SERVER_MODAL_YAML = "modal_server.yaml"

func (srv *DeployServer) ViewIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//build context
	context := make(map[string]interface{})
	context["wizard"] = srv.getInstallWizard()
	context["modal"] = srv.getServerEditModal()

	srv.RespRenderFile(w, "index.html", context)
}

func (srv *DeployServer) getInstallWizard() map[string]string {
	var wizard = make(map[string]string)
	var inlineJsBuffer bytes.Buffer
	var externalJsBuffer bytes.Buffer

	//only iterate steps by name then we can make sure inline js
	//are assembled by order
	for stepName, yamlFilename := range INSTALL_WIZARD {
		yamlContent, err := srv.ReadAssetBytes(yamlFilename)
		if err != nil {
			log.Println("Error read yaml content:", yamlFilename)
			continue
		}

		if ok := srv.y2h.ReadBytes(yamlContent); !ok{
			log.Println("Error parse yaml content:", yamlFilename)
			continue
		}

		//get HTML output
		wizard[stepName] = srv.y2h.GetHtml()

		//get Javascript output
		for jsType, jsContent := range srv.y2h.GetJavascript() {
			switch jsType {
			case "inline":
				inlineJsBuffer.WriteString(jsContent)
			case "external":
				externalJsBuffer.WriteString(jsContent)
			}
		}
	}

	wizard["inlineJs"] = inlineJsBuffer.String()
	wizard["externalJs"] = externalJsBuffer.String()

	return wizard
}

func (srv *DeployServer) getServerEditModal() map[string]string {
	yamlContent, err := srv.ReadAssetBytes(DEPLOY_SERVER_MODAL_YAML)
	if err != nil {
		log.Println("Error read yaml content:", DEPLOY_SERVER_MODAL_YAML)
		return nil
	}

	if ok := srv.y2h.ReadBytes(yamlContent); !ok{
		log.Println("Error parse yaml content:", DEPLOY_SERVER_MODAL_YAML)
		return nil
	}

	serverEditModal := make(map[string]string)
	serverEditModal["html"] = srv.y2h.GetHtml()
	inlineJs, ok := srv.y2h.GetJavascript()["inline"]
	if ok {
		serverEditModal["inlineJs"] = inlineJs
	}

	return serverEditModal
}

