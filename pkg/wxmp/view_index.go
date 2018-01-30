package wxmp

import (
	"net/http"
	"log"
	//"io/ioutil"
	"fmt"
	"github.com/rfancn/wegigo/sdk/wxmp"
)

func (srv *WxmpServer) ViewVerifyWxmpServer(w http.ResponseWriter, r *http.Request) {
	log.Println("verify wxmp server")

	r.ParseForm()
	fmt.Fprintf(w, r.FormValue("echostr"))
}

func (srv *WxmpServer) ViewMain(w http.ResponseWriter, r *http.Request) {
	log.Println("server index")

	wxmpRequest, ok := r.Context().Value("wxmpRequest").(*wxmp.WxmpRequest)
	if !ok {
		log.Println("Invalid wxmp")
		srv.RespText(w, "")
		return
	}

	log.Println(wxmpRequest.MsgId)

	/**
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}
	**/

	//send message to corresponding queue
	//srv.rmqManager.TopicPublish(srv.Name, "AutoReply", "text/plain", body)

	//wait for the result

	srv.RespText(w, "done")
}
