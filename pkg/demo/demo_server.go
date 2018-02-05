package demo

import (
	"log"
	"github.com/rfancn/wegigo/sdk/server"
	//"path"
	"github.com/justinas/alice"
	"github.com/coreos/etcd/clientv3"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type DemoServer struct {
	server.SimpleServer
	remoteConfig  *WegigoRemoteConfig
	etcd     *clientv3.Client
	rabbitmq *ProxyRabbitmq
	stopChan chan struct{}
}

func NewDemoServer(serverName string, assetDir string, etcdAddress string, etcdPort int, rabbitmqAddress string, rabbitmqPort int) *DemoServer {
	wegigoServer := &DemoServer{}

	if ! wegigoServer.SimpleServer.Initialize(serverName, assetDir, Asset, AssetDir, AssetInfo) {
		return nil
	}

	wegigoServer.rabbitmq = NewRabbitmq(rabbitmqAddress, rabbitmqPort)
	if wegigoServer.rabbitmq == nil {
		return nil
	}

	wegigoServer.etcd = NewEtcdClient(etcdAddress, etcdPort)
	if wegigoServer.etcd == nil {
		return nil
	}

	wegigoServer.remoteConfig = NewRemoteConfig(wegigoServer.etcd)

	wegigoServer.stopChan = make(chan struct{})

	return wegigoServer
}


func Run(assetDir string, etcdAddress string, etcdPort int, rabbitmqAddress string, rabbitmqPort int) {
	log.Printf("Run Wegigo Server")

	srv := NewWegigoServer("wxmp", assetDir, etcdAddress, etcdPort, rabbitmqAddress, rabbitmqPort)
	if srv == nil {
		log.Fatal("Error create wegigo server")
	}
	defer srv.Close()

	//setup graceful shutdown handler
	srv.setupShutdownHandler()

	//watch app config changes
	go srv.remoteConfig.WatchApps(srv.etcd, srv.stopChan)

	srv.setupRouter()

	err := srv.RunHttps("0.0.0.0", 443)
	if err != nil {
		log.Fatal("Error start proxy server:", err)
	}

	//log.Println("Enabled Ids:", srv.remoteConfig.GetAppIds("enabled"))
}

func (srv *DemoServer) setupRouter() {
	log.Println("Setup router")

	srv.AddRoute("get", "/", srv.ViewIndex)

	//chain handler: isEnabledMiddleware->mainHandler
	handlerChain := alice.New(isEnabledMiddleware).Then(http.HandlerFunc(mainHandler))
	srv.AddRoute("get", "/app/:Uuid", HandleFunc(handlerChain, srv.remoteConfig.enabledApps))

	//app config route
	srv.AddRoute("get", "/config/app/:Uuid", srv.ViewAppConfig)
}

func (srv *DemoServer) setupShutdownHandler() {
	gracefulShutdownChan := make(chan os.Signal, 2)
	signal.Notify(gracefulShutdownChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-gracefulShutdownChan
		srv.Close()
		os.Exit(1)
	}()
}

func (srv *DemoServer) Close() {
	close(srv.stopChan)
	srv.rabbitmq.Close()
	srv.etcd.Close()
}

