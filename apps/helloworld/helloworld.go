package main

import (
	"log"
	"github.com/rfancn/wegigo/sdk/app"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

var App appHelloWorld
var APP_INFO = &app.AppInfo{
	Uuid: "1234d61b-e9fe-4346-b64d-dcd1102f1234",
	Name: "HelloWorld",
	Version: "0.0.1",
	Author: "Ryan Fan",
	Desc: "echo back",
	Configurable: false,
}

type appHelloWorld struct {
	app.BaseApp
}


func (a *appHelloWorld) 	Init(serverName string, rootDir string, etcdUrl string, amqpUrl string) error {
	return a.BaseApp.Initialize(serverName, rootDir, etcdUrl, amqpUrl, APP_INFO, a)
}

func (a *appHelloWorld) LoadConfig() {
	configData := a.GetConfigData()
	if configData == nil {
		log.Println("[WARN] Failed to get HelloWorld config data from DB!")
		return
	}
}

func (a *appHelloWorld) Match(data []byte) bool{
	return true
}

func (a *appHelloWorld) Process(data []byte) []byte{
	log.Println("HelloWorld received:", string(data))

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/test")

	// if there is an error opening the connection, handle it
	if err != nil {
		log.Println("Error open db", err.Error())
		return data
	}

	// defer the close till after the main function has finished
	// executing
	defer db.Close()


	// perform a db.Query insert
	insert, err := db.Query("insert into equipment (type, color, working, location) values ('ddd', 'red', 1, 'london')")

	// if there is an error inserting, handle it
	if err != nil {
		log.Println("Error insert to db", err.Error())
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()

	return data
}









