package wxmp

//InitEnv: sync server's running env info to etcd which will be read by app plugin
func (srv *WxmpServer) InitEnv() {
	env := &WxmpServerEnv{}
	env.serverName = srv.Name
	env.amqpUrl = srv.cmdArg.RabbitmqUrl
	env.etcdUrl = srv.cmdArg.EtcdUrl
}
