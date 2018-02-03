package wxmp

import (
	"fmt"
	"time"
	"strings"
)

const (
	// Reply format
	tplReplyHeader    = "<ToUserName><![CDATA[%s]]></ToUserName><FromUserName><![CDATA[%s]]></FromUserName><CreateTime>%d</CreateTime>"
	tplReplyText       = "<xml>%s<MsgType><![CDATA[text]]></MsgType><Content><![CDATA[%s]]></Content></xml>"
	tplReplyImage      = "<xml>%s<MsgType><![CDATA[image]]></MsgType><Image><MediaId><![CDATA[%s]]></MediaId></Image></xml>"
	tplReplyVoice      = "<xml>%s<MsgType><![CDATA[voice]]></MsgType><Voice><MediaId><![CDATA[%s]]></MediaId></Voice></xml>"
	tplReplyVideo      = "<xml>%s<MsgType><![CDATA[video]]></MsgType><Video><MediaId><![CDATA[%s]]></MediaId><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description></Video></xml>"
	tplReplyMusic      = "<xml>%s<MsgType><![CDATA[music]]></MsgType><Music><Title><![CDATA[%s]]></Title><Description><![CDATA[%s]]></Description><MusicUrl><![CDATA[%s]]></MusicUrl><HQMusicUrl><![CDATA[%s]]></HQMusicUrl><ThumbMediaId><![CDATA[%s]]></ThumbMediaId></Music></xml>"
	tplReplyNews       = "<xml>%s<MsgType><![CDATA[news]]></MsgType><ArticleCount>%d</ArticleCount><Articles>%s</Articles></xml>"
	tplReplyArticle    = "<item><Title><![CDATA[%s]]></Title> <Description><![CDATA[%s]]></Description><PicUrl><![CDATA[%s]]></PicUrl><Url><![CDATA[%s]]></Url></item>"
	tplTransferCustomerService = "<xml>" + tplReplyHeader + "<MsgType><![CDATA[transfer_customer_service]]></MsgType></xml>"
)

type Reply struct {
	fromUserName string
	header string
}

func NewReply(request *Request) *Reply{
	xmlReplyHeader :=  fmt.Sprintf(tplReplyHeader, request.FromUserName, request.ToUserName, time.Now().Unix())
	return &Reply{fromUserName: request.ToUserName, header: xmlReplyHeader}
}

// Format reply message header
func (r Reply) getReplyHeader(request *Request) string {
	return fmt.Sprintf(tplReplyHeader, request.FromUserName, request.ToUserName, time.Now().Unix())
}

// Reply empty message
func (r *Reply) ReplyOK() []byte{
	return []byte("success")
}

// Reply text message
func (r *Reply) ReplyText(text string) []byte {
	msg := fmt.Sprintf(tplReplyText, r.header, text)
	return []byte(msg)
}

// Reply image message
func (r *Reply) ReplyImage(mediaId string) []byte{
	msg := fmt.Sprintf(tplReplyImage, r.header, mediaId)
	return []byte(msg)
}

// Reply voice message
func (r *Reply) ReplyVoice(mediaId string) []byte{
	msg := fmt.Sprintf(tplReplyVoice, r.header, mediaId)
	return []byte(msg)
}

// Reply video message
func (r *Reply) ReplyVideo(mediaId string, title string, description string) []byte{
	msg := fmt.Sprintf(tplReplyVideo, r.header, mediaId, title, description)
	return []byte(msg)
}

// Reply music message
func (r *Reply) ReplyMusic(m *Music) []byte{
	msg := fmt.Sprintf(tplReplyMusic, r.header, m.Title, m.Description, m.MusicUrl, m.HQMusicUrl, m.ThumbMediaId)
	return []byte(msg)
}

// Reply news message (max 10 news)
func (r *Reply) ReplyNews(articles []Article) []byte{
	articleList := make([]string, 0)
	for _, article := range articles {
		articleList = append(articleList, fmt.Sprintf(tplReplyArticle, article.Title, article.Description, article.PicUrl, article.Url))
	}

	msg := fmt.Sprintf(tplReplyNews, r.header, len(articleList), strings.Join(articleList, ""))
	return []byte(msg)
}

// Transfer customer service
func (r *Reply) TransferCustomerService(serviceId string) []byte {
	msg := fmt.Sprintf(tplTransferCustomerService, serviceId, r.fromUserName, time.Now().Unix())
	return []byte(msg)
}
