package main

import (
	"time"
	"davizzard/ErasureCodes/src/goObjStore/src/httpGo"
	"davizzard/ErasureCodes/src/goObjStore/src/conf"
	"net/http"
)
func main() {

	router := httpGo.MyNewRouter()

	go func(){http.ListenAndServe(conf.Peer2a, router)}()
	go func(){http.ListenAndServe(conf.Peer2b, router)}()
	http.ListenAndServe(conf.Peer2c, router)
	time.Sleep(1*time.Hour)

}

