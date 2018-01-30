package wxmp

import (
	"io/ioutil"
	"log"
	"net/http"
	"encoding/xml"
)

// Weixin MP request
type WxmpRequest struct {
	MessageHeader
	MsgId        int64
	Content      string
	PicUrl       string
	MediaId      string
	Format       string
	ThumbMediaId string
	LocationX    float32 `xml:"Location_X"`
	LocationY    float32 `xml:"Location_Y"`
	Scale        float32
	Label        string
	Title        string
	Description  string
	Url          string
	Event        string
	EventKey     string
	Ticket       string
	Latitude     float32
	Longitude    float32
	Precision    float32
	Recognition  string
	Status       string
}

//NewWxmpRequest: get WxmpRequest from http.Request
func NewWxmpRequest(r *http.Request) *WxmpRequest {
	// Process message
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("NewWxmpRequest(): Error read http request body", err)
		return nil
	}

	req := &WxmpRequest{}
	if err := xml.Unmarshal(data, &req); err != nil {
		log.Println("NewWxmpRequest(): Error unmarshal wxmp message:", err)
		return nil
	}

	return req
}
