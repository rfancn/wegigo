package wxmp

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"github.com/rfancn/goy2h"
	"github.com/rfancn/wegigo/sdk/server"
	"github.com/rfancn/wegigo/sdk/app"
	"github.com/rfancn/wegigo/sdk/rabbitmq"
)

const SERVER_NAME = "wxmp"

//command argument passed from cmd package
type WxmpCmdArgument struct {
	ServerUrl   string
	EtcdUrl 	string
	RabbitmqUrl	string

	//Asset dir
	AssetDir        string
	//application plugins dir
	AppPluginDir 	string
	//test concurrency for app worker, it will moved to etcd config in the future
	AppConcurrency     int
}

type WxmpServerEnv struct {
	serverName string
	etcdUrl string
	amqpUrl string
}

type WxmpServer struct {
	server.SimpleServer
	//command argument
	cmdArg 			*WxmpCmdArgument
	// all kinds of managers
	appManager 		*app.AppManager
	rmqManager 		*rabbitmq.RabbitMQManager

	//store all initialized App instance
	apps        map[string]app.IApp
	appInfos    map[string]*app.AppInfo
	//store all enabled App intance
	enabledApps map[string]app.IApp

	//stop watch channel, all watch routine need check this for quit or not
	stopWatch   chan struct{}

	//yaml to html translator, which will be used in app config view
	y2h *goy2h.Y2H
}

func NewWxmpServer(serverName string, arg *WxmpCmdArgument) *WxmpServer {
	srv := &WxmpServer{}

	if ! srv.SimpleServer.Initialize(serverName, arg.AssetDir, Asset, AssetDir, AssetInfo) {
		return nil
	}

	appManager, err := app.NewAppManager(arg.EtcdUrl)
	if err != nil {
		log.Println("NewWxmpServer(): Error create app manager:", err)
		return nil
	}

	rmqManager, err := rabbitmq.NewRabbitMQManager(arg.RabbitmqUrl)
	if err != nil {
		log.Println("NewWxmpServer(): Error create rabbitmq manager:", err)
		return nil
	}

	//assign to server instance
	srv.cmdArg = arg
	srv.appManager = appManager
	srv.rmqManager = rmqManager
	srv.y2h = goy2h.New()

	//init other varaibles
	srv.apps = make(map[string]app.IApp)
	srv.appInfos = make(map[string]*app.AppInfo)
	srv.enabledApps = make(map[string]app.IApp)

	srv.stopWatch = make(chan struct{})

	return srv
}

func Run(cmdArg *WxmpCmdArgument) {
	log.Printf("Run Wechat Media Platform Server")

	srv := NewWxmpServer(SERVER_NAME, cmdArg)
	if srv == nil {
		log.Fatal("Error create wxmp server")
	}
	defer srv.Close()

	//srv.InitEnv()

	//setup graceful shutdown handler
	srv.setupShutdownHandler()

	//server's amqp must be setup before discovering apps
	srv.SetupAMQP()

	srv.LoadAndRunApps()

	go srv.WatchEnabledApps()
	go srv.WatchAppConfig()

	srv.SetupRouter()

	err := srv.Run(cmdArg.ServerUrl)
	if err != nil {
		log.Fatal("Error start wxmp server:", err)
	}
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
	close(srv.stopWatch)
	srv.rmqManager.Close()
	srv.appManager.Close()
}
