package wxmp

import (
	"net/http"
	"log"
	"time"
)

//ViewVerifyWxmpServer: verify the changes made in Wxmp Platform
func (srv *WxmpServer) ViewVerifyWxmpServer(w http.ResponseWriter, r *http.Request) {
	log.Println("verify wxmp server")

	//echo back "echostr" when succeed
	w.Write([]byte(r.FormValue("echostr")))
}

func (srv *WxmpServer) ViewMain(w http.ResponseWriter, r *http.Request) {
	//1. declare reply queue
	replyQueue := srv.rmqManager.DeclareTempQueue()

	//2. consume
	ch, msgs, err := srv.rmqManager.Consume(replyQueue)
	defer ch.Close()
	if err != nil {
		log.Println("ViewMain(): Error consume wxmp reply queue:", err)
		w.Write([]byte("success"))
		return
	}

	//send message
	nonce := r.Context().Value("nonce").(string)
	data := r.Context().Value("data").([]byte)
	msgHeaders := r.Context().Value("msgHeaders").(map[string]interface{})

	log.Println(msgHeaders)

	ok := srv.rmqManager.RPCPublishJsonWithHeaders(
		srv.Name,
		msgHeaders,
		replyQueue,
		//use nonce as msg's CorrelationId
		nonce,
		data)

	if !ok {
		log.Println(w, "ViewMain(): Error publish received wxmp request")
		w.Write([]byte("success"))
		return
	}

	//monitor reply
	select {
	case m := <-msgs:
		log.Println("ViewMain(): Receive reply:", string(m.Body))
		if m.CorrelationId !=  nonce{
			log.Println("Reply correlationId umatched!")
			w.Write([]byte("success"))
			return
		}
		w.Write(m.Body)
		//wxmp timeout waiting for reply is 5 seconds,
		//but as queue declare/consume/publish may consume some time
		//here use 100 milliseconds compensation
	case <-time.After(4 * time.Second + 900 * time.Millisecond):
		log.Printf("ViewMain(): Timeout waiting for reply: %v", time.Now())
		w.Write([]byte("success"))
	}

}
