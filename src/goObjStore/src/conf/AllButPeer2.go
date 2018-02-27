package conf
import (
"os"
)
// COMPULSORY ----------------------------------------------------------------------------
var NumNodes int = 3
var LocalFileName = "NEW.xml"
var ChunkProxyName = "NEW"
var PortsPerNode int = 3
//const differentIPs = 2
const IP = "192.168.1.10"
var FilePath = os.Getenv("GOPATH")+"/src/davizzard/ErasureCodes/src/goObjStore/src/bigFile"
//var FilePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/dataset.xml"
// TRACKER
var TrackerAddr = IP+":8000"

var Proxy1 = IP+":8070"
var Proxy2 = IP+":8071"
var Proxy3 = IP+":8072"


var ProxyAddr=[]string{Proxy1,Proxy2,Proxy3/*,Proxy4,Proxy5*/}
var DownloadsDirectory=os.Getenv("GOPATH")+"/src/davizzard/ErasureCodes/src/goObjStore/src/Downloads/"
var LocalDirectory=os.Getenv("GOPATH")+"/src/davizzard/ErasureCodes/src/goObjStore/src/local/"
// ---------------------------------------------------------------------------------------


// IMPORTANT NOTE: last character of port is the peer's internal identifier


// THIS MACHINE --------------------------------------------------------------------------
// In order to accept request from all machines, the machine running the Tracker
// needs to have all the I.P. and Addresses in the variable Peers
var Peer1a = IP+":8011"
var Peer1b = IP+":8021"
var Peer1c = IP+":8031"
var Peer1List = []string{Peer1a, Peer1b, Peer1c}


var Peer3a = IP+":8013"
var Peer3b = IP+":8023"
var Peer3c = IP+":8033"
var Peer3List = []string{Peer3a, Peer3b, Peer3c}

var Peers =[][]string{Peer1List, Peer2List, Peer3List/*, Peer4List, Peer5List*/}
// ---------------------------------------------------------------------------------------




// OTHER MACHINES ------------------------------------------------------------------------
const otherIP="192.168.1.11"

var Peer2a = otherIP+":8012"
var Peer2b = otherIP+":8022"
var Peer2c = otherIP+":8032"
var Peer2List = []string{Peer2a, Peer2b, Peer2c}
// ---------------------------------------------------------------------------------------
