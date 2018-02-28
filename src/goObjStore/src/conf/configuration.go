package conf
import (
	"os"
)

var NumNodes int = 3
var LocalFileName = "NEW.xml"
var ChunkProxyName = "NEW"
var PortsPerNode int = 3
const IP = "127.0.0.1"

var FilePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/bigFile"
//var FilePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/dataset.xml"
// TRACKER
var TrackerAddr = IP+":8000"
var DownloadsDirectory=os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/Downloads/"
var LocalDirectory=os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/local/"
// PEERS
// IMPORTANT NOTE: last character of port is the peer's internal identifier
var Peer1a = IP+":8011"
var Peer1b = IP+":8021"
var Peer1c = IP+":8031"
var Peer1List = []string{Peer1a, Peer1b, Peer1c}

var Peer2a = IP+":8012"
var Peer2b = IP+":8022"
var Peer2c = IP+":8032"
var Peer2List = []string{Peer2a, Peer2b, Peer2c}

var Peer3a = IP+":8013"
var Peer3b = IP+":8023"
var Peer3c = IP+":8033"
var Peer3List = []string{Peer3a, Peer3b, Peer3c}
/*
	var Peer4a = IP+":8011"
	var Peer4b = IP+":8011"
	var Peer4c = IP+":8011"
	var Peer4List = []string{Peer4a, Peer4b, Peer4c}

	var Peer5a = IP+":8011"
	var Peer5b = IP+":8011"
	var Peer5c = IP+":8011"
	var Peer5List = []string{Peer5a, Peer5b, Peer5c}
*/
var Peers =[][]string{Peer1List, Peer2List, Peer3List/*, Peer4List, Peer5List*/}

// Proxy
var Proxy1 = IP+":8070"
var Proxy2 = IP+":8071"
var Proxy3 = IP+":8072"
/*
var Proxy4 = IP+":8073"
var Proxy5 = IP+":8074"
*/

var ProxyAddr=[]string{Proxy1,Proxy2,Proxy3/*,Proxy4,Proxy5*/}

