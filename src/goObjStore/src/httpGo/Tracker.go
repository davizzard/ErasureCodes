package httpGo
import(
	"fmt"
	"github.com/davizzard/ErasureCodes/src/goObjStore/src/httpVar"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"io"
	"strconv"
	"strings"
	"math/rand"
)

type NodeNum struct {
	Quantity string		`json:"Quantity"`
	ID string		`json:"ID"`
	Type string 		`json:"Type"`
}
type jsonKey struct {
	ID string		`json:"ID"`
	Type string 		`json:"Type"`
}

/*
StartTracker is called when a Tracker is being initialized.
Registers the nodes' addresses
@param1 List of nodes lists of addresses (each node can have multiple addresses)
	check httpGo/Main.go for details
 */
func StartTracker(nodeList [][]string){
	var nodeAux httpVar.NodeInfo
	for _, node := range nodeList {
		nodeAux.Url=node
		nodeAux.Busy=false
		httpVar.TrackerNodeList = append(httpVar.TrackerNodeList, nodeAux)
	}


}

/*
GetNodes is called when a GET requests [TrackerURL]/GetNodes.
Sends new json encoded node list back to the sender
@param1 used by an HTTP handler to construct an HTTP response.
@param2 represents HTTP request
 */
func GetNodes(w http.ResponseWriter, r *http.Request){
	var nodeNum NodeNum
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	CheckSimpleErr(err, nil, true)
	err = r.Body.Close()
	CheckSimpleErr(err, nil, true)

	err = json.Unmarshal(body, &nodeNum)
	CheckJsonErr(err, nil, w)

	num, err := strconv.Atoi(nodeNum.Quantity)
	CheckJsonErr(err, nil, w)

	nodeList:=chooseNodes(num)
	// registering nodeList to type and ID can be: object, container or account
	if strings.Compare(nodeNum.Type,"object")==0 {
		fmt.Println("object")
		httpVar.MapKeys.Lock()
		httpVar.MapKeyNodes[nodeNum.ID] = nodeList
		httpVar.MapKeys.Unlock()
	}else if strings.Compare(nodeNum.Type,"account")==0{
		httpVar.MapAcc.Lock()
		httpVar.MapAccNodes[nodeNum.ID] = nodeList
		httpVar.MapAcc.Unlock()
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(nodeList); err != nil {
		fmt.Println("GetNodes: error encoding response: ",err.Error())
	}
}

func chooseNodes(num int)[][]string{
	i:=0
	var response [][]string
	var busies [][]string
	for i<num{
		if httpVar.TrackerNodeList[i].Busy==false {
			response = append(response, httpVar.TrackerNodeList[i].Url )
			httpVar.TrackerNodeList[i].Busy=true
		}else{
			busies=append(busies, httpVar.TrackerNodeList[i].Url)
		}
		i++

	}
	if len(response)<num{
		fmt.Println("There is not enough free nodes, adding bussies")
		busyResponse:= chooseBusyNodes(num,busies, response)
		return busyResponse
	}
	return response
}

func chooseBusyNodes(num int, busies [][]string, response [][]string) [][]string{
	var random int // iterates through bussies
	var busiesAssigned = make(map[int]bool) // registers random indexes used
	for len(response)<num && len(busiesAssigned)<len(busies){
		random=rand.Intn(len(busies))
		if len(busiesAssigned)!=0{
			// check if index used
			_, exists := busiesAssigned[random]
			if !exists{
				response=append(response, busies[random])
				busiesAssigned[random]=true
			}
		} else{
			response=append(response, busies[random])
			busiesAssigned[random]=true
		}

	}
	if num>len(response){fmt.Println(" total number of nodes less than proxy asked ")}
		// proxy will receive less nodes than expected
	return response


}

func GetNodesForKey(w http.ResponseWriter, r *http.Request){
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	CheckSimpleErr(err, nil, true)
	err = r.Body.Close()
	CheckSimpleErr(err, nil, true)

	var request jsonKey
	err = json.Unmarshal(body, &request)
	CheckJsonErr(err, nil, w)

	var nodeList [][]string
	if strings.Compare(request.Type, "object") == 0  {
		nodeList = httpVar.MapKeyNodes[request.ID]
		if len(nodeList)==0{
			fmt.Println("ID ",request.ID, " has length 0")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}else if strings.Compare(request.Type, "account" )== 0 {
		nodeList = httpVar.MapAccNodes[request.ID]
		if len(nodeList)==0{
			fmt.Println("ID ",request.ID, " has length 0")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(nodeList); err != nil {
		fmt.Println("GetNodesForKey: error encoding response: ",err.Error())
	}
}




