package wxmp

func (srv *WxmpServer) setupProducer() {
	srv.rmqManager.DeclareTopicExchange(SERVER_NAME, false)
}
