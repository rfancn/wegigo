package wxmp

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"context"
	"github.com/rfancn/wegigo/sdk/server"
	"github.com/rfancn/wegigo/sdk/etcd"
	"github.com/rfancn/wegigo/sdk/app"
	"github.com/rfancn/wegigo/sdk/rabbitmq"
)

const SERVER_NAME = "wxmp"

type WxmpServer struct {
	server.SimpleServer

	// all kinds of managers
	etcdManager 	*etcd.EtcdManager
	appManager 		*app.AppManager
	rmqManager 		*rabbitmq.RabbitMQManager

	appsDir string
	ctx context.Context

	//stop channel to indicate
	stopChan chan struct{}
}

func NewWxmpServer(serverName string, appsDir string, assetDir string, etcdAddress string, etcdPort int, rabbitmqAddress string, rabbitmqPort int) *WxmpServer {
	srv := &WxmpServer{}

	if ! srv.SimpleServer.Initialize(serverName, assetDir, Asset, AssetDir, AssetInfo) {
		return nil
	}

	//new all kinds of manager
	etcdManager := etcd.NewEtcdManager(etcdAddress, etcdPort)
	if etcdManager == nil {
		log.Println("NewWxmpServer(): Error create etcd manager")
		return nil
	}

	appManager := app.NewAppManager(etcdManager)
	if appManager == nil {
		log.Println("NewWxmpServer(): Error create app manager")
		return nil
	}

	rmqManager := rabbitmq.NewRabbitMQManager(rabbitmqAddress, rabbitmqPort)
	if rmqManager == nil {
		log.Println("NewWxmpServer(): Error create rabbitmq manager")
		return nil
	}

	//assign to server instance
	srv.appsDir = appsDir
	srv.stopChan = make(chan struct{})
	srv.etcdManager = etcdManager
	srv.appManager = appManager
	srv.rmqManager = rmqManager

	return srv
}


func Run(appsDir string, assetDir string, etcdAddress string, etcdPort int, rabbitmqAddress string, rabbitmqPort int) {
	log.Printf("Run Wechat Media Platform Server")

	srv := NewWxmpServer(SERVER_NAME, appsDir, assetDir, etcdAddress, etcdPort, rabbitmqAddress, rabbitmqPort)
	if srv == nil {
		log.Fatal("Error create wxmp server")
	}
	defer srv.Close()

	//setup graceful shutdown handler
	srv.setupShutdownHandler()

	srv.discoverApps()

	srv.setupRouter()

	err := srv.RunHttps("0.0.0.0", 443)
	if err != nil {
		log.Fatal("Error start wxmp server:", err)
	}

	//log.Println("Enabled UUIDs:", srv.remoteConfig.GetAppUuids("enabled"))
}

func (srv *WxmpServer) setupRouter() {
	log.Println("Setup router")

	srv.AddRoute("get", "/admin/", srv.ViewIndex)
	srv.AddRoute("get", "/admin/app/", srv.ViewAppIndex)

	//chain handler: isEnabledMiddleware->mainHandler
	//handlerChain := alice.New(isEnabledMiddleware).Then(http.HandlerFunc(mainHandler))
	//srv.AddRoute("get", "/app/:uuid", HandleFunc(handlerChain, srv.remoteConfig.enabledApps))

	//app config route
	srv.AddRoute("get", "/admin/app/config/:uuid", srv.ViewAppConfig)
	srv.AddRoute("post", "/admin/app/toggle/:uuid", srv.ViewAppToggle)

}

func (srv *WxmpServer) setupShutdownHandler() {
	gracefulShutdownChan := make(chan os.Signal, 2)
	signal.Notify(gracefulShutdownChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-gracefulShutdownChan
		srv.Close()
		os.Exit(1)
	}()
}

func (srv *WxmpServer) Close() {
	close(srv.stopChan)
	srv.rmqManager.Close()
	srv.etcdManager.Close()
}
