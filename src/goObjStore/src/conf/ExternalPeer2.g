package conf
import (
"os"
)
// COMPULSORY --------------------------------------------------------------------------------------
var NumNodes int = 3
var LocalFileName = "NEW.xml"
var ChunkProxyName = "NEW"
var PortsPerNode int = 3
const IP = "192.168.1.11"
//const differentIPs=2
var FilePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/bigFile"
//var FilePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/dataset.xml"
// TRACKER
var TrackerAddr = IP+":8000"

var Proxy1 = IP+":8070"
var Proxy2 = IP+":8071"
var Proxy3 = IP+":8072"


var ProxyAddr=[]string{Proxy1,Proxy2,Proxy3/*,Proxy4,Proxy5*/}
var DownloadsDirectory=os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/Downloads/"
// ----------------------------------------------------------------------------------------------------


// IMPORTANT NOTE: last character of port is the peer's internal identifier


// THIS MACHINE ---------------------------------------------------------------------------------------
// In order to accept request from all machines, the machine running the Tracker
// needs to have all the IP Addresses in the variable Peers
var Peer2a = IP+":8012"
var Peer2b = IP+":8022"
var Peer2c = IP+":8032"
// ----------------------------------------------------------------------------------------------------


// OTHER MACHINES -------------------------------------------------------------------------------------
// Please follow the otherIP1, otherIP2, otherIP3 sequence, must be a const too
const otherIP1 = "192.168.1.11"
// ----------------------------------------------------------------------------------------------------
