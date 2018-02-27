package main

import (
	"os"
	"time"
	"davizzard/ErasureCodes/src/goObjStore/src/httpGo"
	"net/http"

)

func main2() {

	const IP = "127.0.0.1"

	//var filePath2 = os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/bigFile"
	var filePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/dataset.xml"


	var proxyAddr = "127.0.0.1:8080"
	var trackerAddr = IP+":8000"

	// PEERS
	// IMPORTANT NOTE: last character of port is the peer's internal identifier
	var Peer1a = IP+":8011"
	var Peer1b = IP+":8021"
	var Peer1c = IP+":8031"
	var peer1List = []string{Peer1a, Peer1b, Peer1c}

	var Peer2a = IP+":8012"
	var Peer2b = IP+":8022"
	var Peer2c = IP+":8032"
	var peer2List = []string{Peer2a, Peer2b, Peer2c}

	var Peer3a = IP+":8013"
	var Peer3b = IP+":8023"
	var Peer3c = IP+":8033"
	var peer3List = []string{Peer3a, Peer3b, Peer3c}
	/*
		var Peer4a = IP+":8011"
		var Peer4b = IP+":8011"
		var Peer4c = IP+":8011"
		var peer4List = []string{Peer4a, Peer4b, Peer4c}

		var Peer5a = IP+":8011"
		var Peer5b = IP+":8011"
		var Peer5c = IP+":8011"
		var peer5List = []string{Peer5a, Peer5b, Peer5c}
	*/
	var Peers =[][]string{peer1List, peer2List, peer3List/*, peer4List, peer5List*/}

	// PROXY
	var proxy1 = IP+":8070"
	var proxy2 = IP+":8071"
	var proxy3 = IP+":8072"
	/*
	var proxy4 = IP+":8073"
	var proxy5 = IP+":8074"
	*/

	var ProxyAddr=[]string{proxy1,proxy2,proxy3/*,proxy4,proxy5*/}

	httpGo.StartTracker(Peers)



	routerTracker := httpGo.MyNewRouter()
	routerPeer := httpGo.MyNewRouter()
	go func(){httpGo.PutNoP2P(filePath, proxyAddr, trackerAddr, 3)
		time.Sleep(5*time.Second)

		httpGo.GetNoP2P("0527cbea2805d89c6d5d6457b7f9f77c",ProxyAddr, trackerAddr)

	}()
	go func(){http.ListenAndServe(Peer1a, routerPeer)}()
	go func(){http.ListenAndServe(Peer1b, routerPeer)}()
	go func(){http.ListenAndServe(Peer1c, routerPeer)}()

	go func(){http.ListenAndServe(Peer2a, routerPeer)}()
	go func(){http.ListenAndServe(Peer2b, routerPeer)}()
	go func(){http.ListenAndServe(Peer2c, routerPeer)}()

	go func(){http.ListenAndServe(Peer3a, routerPeer)}()
	go func(){http.ListenAndServe(Peer3b, routerPeer)}()
	go func(){http.ListenAndServe(Peer3c, routerPeer)}()


	go func(){http.ListenAndServe(proxy1, routerPeer)}()
	go func(){http.ListenAndServe(proxy2, routerPeer)}()
	go func(){http.ListenAndServe(proxy3, routerPeer)}()



	http.ListenAndServe(trackerAddr, routerTracker)

	}

