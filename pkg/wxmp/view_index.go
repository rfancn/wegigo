package wxmp

import (
	"net/http"
	"log"
	"fmt"
	"time"
)

//ViewVerifyWxmpServer: verify the changes made in Wxmp Platform
func (srv *WxmpServer) ViewVerifyWxmpServer(w http.ResponseWriter, r *http.Request) {
	log.Println("verify wxmp server")

	r.ParseForm()

	fmt.Fprintf(w, r.FormValue("echostr"))
}

func (srv *WxmpServer) ViewMain(w http.ResponseWriter, r *http.Request) {
	replyQueueName := srv.rmqManager.DeclareTempQueue()
	ch, msgs, err := srv.rmqManager.Consume(replyQueueName)
	defer ch.Close()
	if err != nil {
		log.Println("ViewMain(): Error consume wxmp reply queue:", err)
		w.Write([]byte("success"))
		return
	}

	//send message
	data := r.Context().Value("data").([]byte)
	msgHeaders := r.Context().Value("msgHeaders").(map[string]interface{})
	ok := srv.rmqManager.RPCPublishJsonWithHeaders(
		srv.Name,
		msgHeaders,
		replyQueueName,
		"123456789",
		data)

	if !ok {
		log.Println(w, "ViewMain(): Error publish received wxmp request")
		w.Write([]byte("success"))
		return
	}


	select {
	case m := <-msgs:
		log.Println("ViewMain(): Receive reply:", string(m.Body))
		if m.CorrelationId != "123456789" {
			log.Println("Reply correlationId id umatched!")
			return
		}
		w.Write(m.Body)
		break
	//timeout after waiting for 5 seconds, as publish may consume some time
	//here use 4.8 seconds
	case <-time.After(4 * time.Second + 800 * time.Millisecond):
		log.Println("ViewMain(): Timeout waiting for reply")
		w.Write([]byte("success"))
		break
	}


}
