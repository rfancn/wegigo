package main

import (
	"net/http"
	"log"
	"path"
)

func (a *autoReplyApp) ViewConfigIndex(w http.ResponseWriter, r *http.Request) {
	log.Println("autoreply index")

	appRoot := r.Context().Value("AppRoot").(string)
	resPath := path.Join(appRoot, "autoreply/asset/html/test.html")

	data, err := Asset(resPath)
	if err != nil {
		log.Println("AutoReplyApp: Error retrieve data:", err)
		return
	}

	w.Write(data)
}
