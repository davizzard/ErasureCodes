package httpVar

import (
	"sync"
)
var CurrentPart = 0
var P2pPart = 0
var TrackerMutex = &sync.Mutex{}
var PeerMutex = &sync.Mutex{}
var DirMutex = &sync.Mutex{}
var GetMutex = &sync.Mutex{}
var TotalNumMutex = &sync.Mutex{} //unnecessary now
var AccFileMutex = &sync.Mutex{}
var AccFileMutexP2P = &sync.Mutex{}


type NodeInfo struct {
	Url []string
	Busy bool
}
var TrackerNodeList []NodeInfo

var MapKeys = &sync.Mutex{}
var MapAcc = &sync.Mutex{}


var MapKeyNodes = make(map[string][][]string)
var MapAccNodes = make(map[string][][]string)

var SendReady = make(chan int, 180)
var SendP2PReady = make(chan int, 20)

var NumGetsMap = make(map[int]int)
