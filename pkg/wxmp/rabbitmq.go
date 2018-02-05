package wxmp

func (srv *WxmpServer) SetupAMQP() {
	//set a topic exchange
	srv.rmqManager.DeclareHeadersExchange(srv.Name)

	//declare a reply queue,it receive all app's reply message
	//qName := srv.rmqManager.DeclareQueue(srv.replyQueueName, false)

	//bind reply queue to exchange with "reply.*" routing key
	//srv.rmqManager.BindQueue(qName, srv.Name, "reply.*")

	/**
	//for enabled apps, declare it's queue
	for Uuid, name := range srv.appManager.GetEnabledApps() {
		log.Println("Declare queue for enabled app:", name)
		srv.rmqManager.DeclareQueue(Uuid, false)
	}
	**/
}
