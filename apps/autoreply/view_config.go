package main

import (
	"net/http"
	"log"
	"path"
)

func (a *autoReplyApp) ViewConfigIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("autoreply index")


	resPath := path.Join(a.Env.RootDir, "asset/html/test.html")

	data, err := Asset(resPath)
	if err != nil {
		log.Println("AutoReplyApp: Error retrieve data:", err)
		return
	}

	w.Write(data)
}
