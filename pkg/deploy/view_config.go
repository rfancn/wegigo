package deploy

import (
	"net/http"
	"log"
	"encoding/json"
	"html/template"
	"os"
)


var HOST_INVENTORY_FILE = "hosts"

type ConfigGeneral struct {
	Debug string `json:"debug"`
}

type ConfigWechat struct {
	Url string 			`json:"url"`
	Token string 		`json:"token"`
	Key string 			`json:"key"`
	Method string 		`json:"method"`
}

type ConfigServer struct {
	Host string 		`json:"host"`
	Port uint 			`json:"port"`
	Username string 	`json:"username"`
	Password string 	`json:"password"`
	Roles []string 		`json:"roles"`
}

type DeployConfig struct {
	General ConfigGeneral 	`json:"general"`
	Wechat  ConfigWechat 	`json:"wechat"`
	Servers []ConfigServer 	`json:"server"`
}


type DeployNodes struct {
	AllNodes    	[]ConfigServer
	MasterNodes 	[]ConfigServer
	WorkerNodes 	[]ConfigServer
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}


func getDeployNodes(servers []ConfigServer) *DeployNodes {
	nodes := &DeployNodes{AllNodes:servers}
	for _, s := range servers {
		if contains(s.Roles, "master") {
			nodes.MasterNodes = append(nodes.MasterNodes, s)
		}

		if contains(s.Roles, "worker") {
			nodes.WorkerNodes= append(nodes.WorkerNodes, s)
		}
	}

	return nodes
}

func (srv *DeployServer) saveConfigInventory(config *DeployConfig) bool {
	log.Println("Save host inventory")

	inventoryTemplate, err := srv.ReadAssetBytes("inventory.tmpl")
	if err != nil {
		log.Printf("Error read inventory template file\n%s", err)
		return false
	}

	f, err := os.Create(HOST_INVENTORY_FILE)
	defer f.Close()
	if err != nil {
		log.Println("Error open file")
		return false
	}

	tmpl, err := template.New("hosts").Parse(string(inventoryTemplate))
	if err != nil {
		log.Printf("Error parse host inventory template string\n%s", err)
		return false
	}

	nodes := getDeployNodes(config.Servers)
	if err = tmpl.Execute(f, nodes); err != nil {
		log.Printf("Error render host inventory template\n%s", err)
		return false
	}

	return true
}

func (srv *DeployServer) ViewConfig(w http.ResponseWriter, r *http.Request) {
	log.Println("Enter ViewConfig, server status is:", srv.status)
	srv.status = STATUS_CONFIG_STARTED

	config := &DeployConfig{}

	//make sure the srv'status assign to the correct value
	defer func(){
		switch srv.response.Result {
		case "success":
			srv.status = STATUS_CONFIG_DONE
		default:
			srv.status = STATUS_INIT
		}
		log.Println("Exit ViewConfig, server status is:", srv.status)
	}()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(config); err != nil {
		log.Println("Error decode deploy config data")
		srv.RespJsonError(w, "Invalid deploy config data")
		return
	}

	if ! srv.saveConfigInventory(config) {
		log.Println("Error save deploy config")
		srv.RespJsonError(w, "Failed to save deploy config")
		return
	}

	srv.RespJsonSuccess(w, "/deploy/deploy")
}
