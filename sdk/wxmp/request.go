package wxmp

import (
	"log"
	"encoding/xml"
)

// Common message header
type RequestHeader struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
}

// Weixin MP request
type Request struct {
	RequestHeader
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

//NewWxmpRequest: get Request from http.Request
func NewRequest(data []byte) *Request {
	req := &Request{}
	if err := xml.Unmarshal(data, &req); err != nil {
		log.Println("NewWxmpRequest(): Error unmarshal wxmp message:", err)
		return nil
	}

	return req
}
