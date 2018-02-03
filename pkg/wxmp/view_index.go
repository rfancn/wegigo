package wxmp

import (
	"net/http"
	"log"
	"fmt"
	"time"
	"encoding/json"
	"github.com/rfancn/wegigo/sdk/wxmp"
	"io/ioutil"
	"strconv"
)

//ViewVerifyWxmpServer: verify the changes made in Wxmp Platform
func (srv *WxmpServer) ViewVerifyWxmpServer(w http.ResponseWriter, r *http.Request) {
	log.Println("verify wxmp server")

	r.ParseForm()

	fmt.Fprintf(w, r.FormValue("echostr"))
}

func (srv *WxmpServer) ViewTest(w http.ResponseWriter, r *http.Request) {
	// Get http body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error read http body ", err)
		w.Write([]byte("success"))
		return
	}

	//new wxmp requst
	wxmpRequest := wxmp.NewRequest(data)
	if wxmpRequest == nil {
		log.Println("Error new WxmpRequest")
		w.Write([]byte("success"))
		return
	}

	msgHeaders := make(map[string]interface{})
	for uuid, app := range srv.appLoader.enabledAppMap {
		msgHeaders[uuid] = strconv.FormatBool(app.Match(data))
	}

	log.Println(msgHeaders)

	stopWaitingReply := make(chan struct{})
	quit := make(chan struct{})

	replyQueueName := srv.rmqManager.DeclareTempQueue()

	//spawn a routine monitor for result
	go func(stop chan struct{}) {
		//wait for the reply
		ch, msgs, err := srv.rmqManager.Consume(replyQueueName)
		if err != nil {
			log.Println("ViewMain(): Error consume wxmp reply queue:", err)
			w.Write([]byte("success"))
			close(quit)
		}
		//close rabbitmq channel, which will destroy consumer
		defer ch.Close()
		//close quit channel to indicate waiting routing exited, and main routint can be quit now too
		defer close(quit)

		select {
		case <-stop:
			log.Println("ViewMain(): Interrupt waiting for reply")
			w.Write([]byte("success"))
			break
		case m := <-msgs:
			log.Println("ViewMain(): Receive reply:", string(m.Body))
			if m.CorrelationId != string(wxmpRequest.MsgId) {
				log.Println("Reply correlationId id umatched!")
				return
			}
			w.Write(m.Body)
			break
			//timeout after waiting for 5 seconds
		case <-time.After(5 * time.Second):
			log.Println("ViewMain(): Timeout waiting for reply")
			w.Write([]byte("success"))
			break
		}

		log.Println("Exit monitor reply")

	}(stopWaitingReply)

	//marshal message to corresponding queue
	amqpMsg, err := json.Marshal(wxmpRequest)
	if err != nil {
		log.Println(w, "ViewMain(): Error marshal received wxmp request:", err)
		close(stopWaitingReply)
		return
	}

	ok := srv.rmqManager.RPCPublishJsonWithHeaders(
		srv.Name,
		msgHeaders,
		replyQueueName,
		string(wxmpRequest.MsgId),
		amqpMsg)

	if !ok {
		log.Println(w, "ViewMain(): Error publish received wxmp request")
		close(stopWaitingReply)
		return
	}

	//wait for reply monitor routine send the signal
	select {
	case <-quit:
		log.Println("ViewMain(): Reply Monitor routine quited, Exit now")
	case <- time.After(10 * time.Second):
		log.Println("ViewMain(): Timeout wait for reply monitor routine, Exit now")
	}

}

//ViewProxy: proxy the wxmp request to app, and respond based on app's reply
//RabbitMQManager:
// - one connection reprents TCP socket connection to broker
// - each channel represents a logic thread which can operate with broker
func (srv *WxmpServer) ViewMain(w http.ResponseWriter, r *http.Request) {
	log.Println("server index")

	// Get http body
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error read http body ", err)
		w.Write([]byte("success"))
		return
	}

	//new wxmp requst
	wxmpRequest := wxmp.NewRequest(data)
	if wxmpRequest == nil {
		log.Println("Error read http body ", err)
		w.Write([]byte("success"))
		return
	}
	log.Println("Received wxmp requets:", wxmpRequest.MsgId, wxmpRequest.MsgType, wxmpRequest.Content)

	stopWaitingReply := make(chan struct{})
	quit := make(chan struct{})

	replyQueueName := srv.rmqManager.DeclareTempQueue()

	//spawn a routine monitor for result
	go func(stop chan struct{}) {
		//wait for the reply
		ch, msgs, err := srv.rmqManager.Consume(replyQueueName)
		if err != nil {
			log.Println("ViewMain(): Error consume wxmp reply queue:", err)
			w.Write([]byte("success"))
			close(quit)
		}
		//close rabbitmq channel, which will destroy consumer
		defer ch.Close()
		//close quit channel to indicate waiting routing exited, and main routint can be quit now too
		defer close(quit)

		select {
		case <-stop:
			log.Println("ViewMain(): Interrupt waiting for reply")
			w.Write([]byte("success"))
			break
		case m := <-msgs:
			log.Println("ViewMain(): Receive reply:", string(m.Body))
			if m.CorrelationId != string(wxmpRequest.MsgId) {
				log.Println("Reply correlationId id umatched!")
				return
			}
			w.Write(m.Body)
			break
		//timeout after waiting for 5 seconds
		case <-time.After(5 * time.Second):
			log.Println("ViewMain(): Timeout waiting for reply")
			w.Write([]byte("success"))
			break
		}

		log.Println("Exit monitor reply")

	}(stopWaitingReply)

	//marshal message to corresponding queue
	amqpMsg, err := json.Marshal(wxmpRequest)
	if err != nil {
		log.Println(w, "ViewMain(): Error marshal received wxmp request:", err)
		close(stopWaitingReply)
		return
	}

	//routine key is: <msg type>.*
	//ok := srv.rmqManager.TopicPublishJson(srv.Name, wxmpRequest.MsgType+".*", amqpMsg)
	headers := map[string]interface{}{
		"type": wxmpRequest.MsgType,
	}

	ok := srv.rmqManager.RPCPublishJsonWithHeaders(
		srv.Name,
		headers,
		replyQueueName,
		string(wxmpRequest.MsgId),
		amqpMsg)

	if !ok {
		log.Println(w, "ViewMain(): Error publish received wxmp request")
		close(stopWaitingReply)
		return
	}

	//wait for reply monitor routine send the signal
	select {
	case <-quit:
		log.Println("ViewMain(): Reply Monitor routine quited, Exit now")
	case <- time.After(10 * time.Second):
		log.Println("ViewMain(): Timeout wait for reply monitor routine, Exit now")
	}

}
