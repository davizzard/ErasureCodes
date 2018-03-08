package httpGo
import(
	"net/http"
	"fmt"
	"io/ioutil"
	"io"
	"log"
	"os"
	"encoding/json"
	"github.com/davizzard/ErasureCodes/src/goObjStore/src/httpVar"
	"strconv"
	"path/filepath"
	"strings"
	"sync"
	"math/rand"
	"github.com/davizzard/ErasureCodes/src/goObjStore/src/conf"
)
var path = (os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore")

func prepSN(w http.ResponseWriter, r *http.Request){
	var hashRecived putObjMSg
	var nodeID int =int(r.Host[len(r.Host)-1]-'0')
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading ",err)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("error body ",err)
	}
	if err := json.Unmarshal(body, &hashRecived); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("error unmarshalling ",err)
		}
	}

	httpVar.DirMutex.Lock()
	// if data directory doesn't exist, create it
	_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/data")
	if err != nil {
		os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/data",0777)
	}

	// if data/chunk.Hash directory doesn't exist, create it
	_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/data/"+hashRecived.ID)
	if err != nil {
		os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/data/"+hashRecived.ID,0777)
	}

	// if data/chunk.Hash/nodeID directory doesn't exist, create it
	_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/data/"+hashRecived.ID +"/"+strconv.Itoa( nodeID))
	if err != nil {
		err2:=os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/data/"+hashRecived.ID +"/"+strconv.Itoa( nodeID),0777)
		if err2!=nil{
			fmt.Println("StorageNode error making dir", err.Error())
		}
	} else{
		os.RemoveAll(os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/data/"+hashRecived.ID +"/"+strconv.Itoa( nodeID))
		err2:=os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/data/"+hashRecived.ID +"/"+strconv.Itoa( nodeID),0777)
		if err2!=nil{
			fmt.Println("StorageNode error making dir", err.Error())
		}
	}
	httpVar.DirMutex.Unlock()
}

func SNPutObj(w http.ResponseWriter, r *http.Request){
	var listenWg sync.WaitGroup
	listenWg.Add(1)
	go func() {
		var chunk msg
		// Get node ID
		var nodeID int = int(r.Host[len(r.Host) - 1] - '0')

		// Listen to tracker
		if r.Method == http.MethodPost {
                        var wg sync.WaitGroup



			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Println("error reading ", err)
			}
			if err := r.Body.Close(); err != nil {
				fmt.Println("error body ", err)
			}
			if err := json.Unmarshal(body, &chunk); err != nil {
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(422) // unprocessable entity
				log.Println(err)
				if err := json.NewEncoder(w).Encode(err); err != nil {
					fmt.Println("error unmarshalling ", err)
				}
			}



			// Save chunk to file
			err = ioutil.WriteFile(path + "/src/data/" + chunk.Hash + "/" + strconv.Itoa(nodeID) + "/NEW" + strconv.Itoa(chunk.Name), chunk.Text, 0777)
			if err != nil {
				fmt.Println("StorageNodeListen: error creating/writing file1", err.Error())
			}
			_, err = os.Stat(path + "/src/data/" + chunk.Hash + "/" + strconv.Itoa(nodeID) + "/NEW" + strconv.Itoa(chunk.Name))
			for err != nil {
				err = ioutil.WriteFile(path + "/src/data/" + chunk.Hash + "/" + strconv.Itoa(nodeID) + "/NEW" + strconv.Itoa(chunk.Name), chunk.Text, 0777)
				if err != nil {
					fmt.Println("SNPutObj: Peer: error creating/writing file p2p2", err.Error())
				}
				fmt.Println("for", chunk.Name)
				_, err = os.Stat(path + "/src/data/" + chunk.Hash + "/" + strconv.Itoa(nodeID) + "/NEW" + strconv.Itoa(chunk.Name))

			}

                        wg.Add(len(chunk.NodeList))

			// Send chunk to peers
			// sending only one chunk to the rest of peers once, don't need to use multiple addr per peer
			var currentAddr int = rand.Intn(len(chunk.NodeList))
			for _, peer := range chunk.NodeList {
				peerURL := "http://" + peer[currentAddr] + "/SNPutObjP2PRequest"

				go func(p string, URL string) {
					if nodeID == int(p[len(p) - 1] - '0') {
						// Don't send to itself

					} else {
						rpipe, wpipe := io.Pipe() // create pipe
						go func() {
							err := json.NewEncoder(wpipe).Encode(&chunk)
							wpipe.Close()                        // close pipe when go routine finishes
							if err != nil {
								fmt.Println("Error encoding to pipe ", err.Error())
							}
						}()
						httpVar.SendP2PReady <- 1
						_, err := http.Post(peerURL, "application/json", rpipe)
						<-httpVar.SendP2PReady
						if err != nil {
							fmt.Println("Error sending http POST p2p", err.Error())
						}
					}

					defer wg.Done()
				}(peer[0], peerURL)
			}
			wg.Wait()
			chunk=msg{}

		}
		listenWg.Done()
	}()
	listenWg.Wait()

}





// Listen to other peers
func SNPutObjP2PRequest(w http.ResponseWriter, r *http.Request) {
	var chunk msg
	// Get peer ID
	var peerID int = int(r.Host[len(r.Host) - 1] - '0')
	// Listen to tracker
	if r.Method == http.MethodPost {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error reading ", err)
		}
		if err := r.Body.Close(); err != nil {
			fmt.Println("error body ", err)
		}
		if err := json.Unmarshal(body, &chunk); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			log.Println(err)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				fmt.Println("error unmarshalling ", err)
			}
		}
		// Save chunk to file
		httpVar.DirMutex.Lock()

		err = ioutil.WriteFile(path + "/src/data/"+chunk.Hash+"/"+strconv.Itoa( peerID)+ "/P2P" + strconv.Itoa(chunk.Name),chunk.Text, 0777)
		if err != nil {
			fmt.Println("SNPutObjP2PRequest: Peer: error creating/writing file p2p3", err.Error())
		}
		_, err = os.Stat(path+ "/src/data/"+chunk.Hash+"/"+strconv.Itoa( peerID)+ "/P2P" + strconv.Itoa(chunk.Name))
		if err != nil {
			err = ioutil.WriteFile(path + "/src/data/"+chunk.Hash+"/"+strconv.Itoa( peerID)+ "/P2P" + strconv.Itoa(chunk.Name), chunk.Text, 0777)
                	if err != nil {
                      		fmt.Println("SNPutObjP2PRequest: Peer: error creating/writing file p2p4", err.Error())
                	}
		}
		httpVar.DirMutex.Unlock()



	}
}

func SNObjGetChunks(w http.ResponseWriter, r *http.Request){
	fmt.Println("GetChunks init.")
	// Get node ID
	var nodeIDint int = int(r.Host[len(r.Host) - 1] - '0')
	var nodeID string = strconv.Itoa(nodeIDint)

	var keyURL jsonKeyURL
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(body, &keyURL); err != nil {
		fmt.Println("GetChunks error Unmashalling ",err.Error())
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(423) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("GetChunks: error unprocessable entity: ",err.Error())
			return
		}
		return
	}


	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	go SNPutObjSendChunksToProxy(nodeID, keyURL.NodeList, keyURL.Key, keyURL.URL, keyURL.NumParts, keyURL.NumParity, keyURL.GetID, keyURL.ShardsInNodes)

	fmt.Println("GetChunks end.")
}

type getMsg struct {
	Text   []byte
	Name   string
	NodeID string
	NodeList [][]string
	Key    string
	Parts  int
	Parity  int
	GetID int
	ShardsInNodes int
}
func SNPutObjSendChunksToProxy(nodeID string, nodeList [][]string, key string, URL string, PartsNum int, ParityNum int, getID int, shardsInNodes int){
	//var wg *sync.WaitGroup
	(*httpVar.WaitingGroupNodes[getID]).Add(1)


	proxyURL:="http://"+URL

	// for each proxy-name in directory send
	filepath.Walk(path + "/src/data/"+key+"/"+nodeID, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Storage Node with ID %q is unavailable: %v\n", nodeID, err)
		} else {
			if strings.Contains(info.Name(),conf.ChunkProxyName) {
				partBuffer:=make([]byte,info.Size())
				file, err := os.Open(path)
				if err != nil {
					fmt.Println("sendChunksToProxy error opening file ",err.Error())
					return nil
				}
				_, err = file.Read(partBuffer)
				if err != nil {
					fmt.Println("sendChunksToProxy error opening file ",err.Error())
					return nil
				}
				m:=getMsg{Text:partBuffer, Name: info.Name(), NodeID:nodeID, NodeList:nodeList, Key:key, Parts:PartsNum, Parity:ParityNum, GetID:getID, ShardsInNodes:shardsInNodes}
				r, w :=io.Pipe()			// create pipe
				go func() {
					defer w.Close()			// close pipe when go routine finishes
					// save buffer to object
					err=json.NewEncoder(w).Encode(&m)
					if err != nil {
						fmt.Println("Error encoding to pipe ", err.Error())
						return
					}
				}()
				res, err := http.Post(proxyURL,"application/json", r )
				if err != nil {
					fmt.Println("sendChunksToProxy: error creating request: ",err.Error())
					return nil
				}
				if err := res.Body.Close(); err != nil {
					fmt.Println(err)
					return nil
				}
			}
		}


		return nil
	})
	(*httpVar.WaitingGroupNodes[getID]).Done()
	fmt.Println("SNPutObjSendChunksToProxy ended.")
}

type MarshalledAcc struct {
	Bytes []byte
	Name string
}

func SNPutAcc(w http.ResponseWriter, r *http.Request){
	fmt.Println("SNPutAcc")
		var accInfo AccInfo
		var peerID int = int(r.Host[len(r.Host) - 1] - '0')
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error reading ", err)
		}
		if err := r.Body.Close(); err != nil {
			fmt.Println("error body ", err)
		}
		if err := json.Unmarshal(body, &accInfo); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			log.Println(err)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				fmt.Println("error unmarshalling ", err)
			}
		}

		// create new account
		conts := make(map[string]Container)
		newAcc:=Account{Name:accInfo.AccName, Containers:conts}

		// marshall new account
		accountBytes, err := newAcc.MarshalMsg(nil)

		marshalledAcc:=MarshalledAcc{Bytes:accountBytes, Name:accInfo.AccName}

		// save account to file
		httpVar.AccFileMutex.Lock()
		err = ioutil.WriteFile(path+"/src/Account"+accInfo.AccName +strconv.Itoa(peerID),accountBytes,0777)

		// send marhshalled account to other nodes to save it
		var wg sync.WaitGroup
		wg.Add(len(accInfo.NodeList))
		var currentAddr int = rand.Intn(len(accInfo.NodeList))
		for _, peer := range accInfo.NodeList{
			peerURL := "http://" + peer[currentAddr] + "/SNPutAccP2PRequest"
			go func(p string, URL string){
				if peerID == int(p[len(p) - 1] - '0') {
					// Don't send to itself
				} else {
					r, w := io.Pipe()
					go func() {
						// save buffer to object
						err = json.NewEncoder(w).Encode(marshalledAcc)
						if err != nil {
							fmt.Println("Error encoding to pipe ", err.Error())
							return
						}
						defer w.Close()                        // close pipe //when go routine finishes
					}()
					response, err := http.Post(peerURL, "application/json", r)
					if err != nil {
						fmt.Println("Error sending http POST p2p", err.Error())
					}
					responseCode:=response.StatusCode
					if responseCode == 201 {
						fmt.Println(response.StatusCode," SN created")
					}
					fmt.Println(responseCode)

				}
				wg.Done()
			}(peer[0], peerURL)
		}
		wg.Wait()
		httpVar.AccFileMutex.Unlock()



}

// Listen to other peers
func SNPutAccP2PRequest(w http.ResponseWriter, r *http.Request) {
	var marshalledAcc MarshalledAcc
	// Get peer ID
	var peerID int = int(r.Host[len(r.Host) - 1] - '0')

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading ", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("error body ", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	if err := json.Unmarshal(body, &marshalledAcc); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("error unmarshalling ", err)
		}
		w.WriteHeader(http.StatusBadRequest)
	}
	// check Accounts map for new Account
	httpVar.AccFileMutexP2P.Lock()
	err = ioutil.WriteFile(path+"/src/Account"+marshalledAcc.Name+strconv.Itoa(peerID),marshalledAcc.Bytes,0777)
	httpVar.AccFileMutexP2P.Unlock()

	if err != nil {
		fmt.Println("SNPutAccP2PRequest: ",err)
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusCreated)

}


func SNGetAcc(w http.ResponseWriter, r *http.Request){

	var accInfo AccInfo
	var peerID int = int(r.Host[len(r.Host) - 1] - '0')
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading ", err)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("error body ", err)
	}
	if err := json.Unmarshal(body, &accInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("error unmarshalling ", err)
		}
		return
	}



	var accPath string = (os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore")

	// read account from file
	httpVar.AccFileMutex.Lock()
	accountBytes, err := ioutil.ReadFile(accPath+"/src/Account"+accInfo.AccName +strconv.Itoa(peerID))
	if err != nil {
		fmt.Println(accPath)
		fmt.Println("SNGetAcc Error: reading ", accPath+"/src/Account"+accInfo.AccName +strconv.Itoa(peerID))
	}
	httpVar.AccFileMutex.Unlock()

	account:=Account{}
	_,err = account.UnmarshalMsg(accountBytes)
	if err != nil {
		fmt.Println("SNGetAcc: error Unmarshalling account")
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(account); err != nil {
		fmt.Println("GetNodes: error encoding response: ",err.Error())
	}



}



func SNPutCont(w http.ResponseWriter, r *http.Request){

		var accInfo AccInfo
		var peerID int = int(r.Host[len(r.Host) - 1] - '0')
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error reading ", err)
		}
		if err := r.Body.Close(); err != nil {
			fmt.Println("error body ", err)
		}
		if err := json.Unmarshal(body, &accInfo); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			log.Println(err)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				fmt.Println("error unmarshalling ", err)
			}
		}





		// read account from file
		httpVar.AccFileMutex.Lock()
		accountBytes, err := ioutil.ReadFile(path+"/src/Account"+accInfo.AccName +strconv.Itoa(peerID))
		if err != nil {
			fmt.Println("SNPutCont Error: reading ", path+"/src/Account"+accInfo.AccName +strconv.Itoa(peerID))
		}

		// update account
		account:=Account{}
		_,err = account.UnmarshalMsg(accountBytes)
		objs:= make(map[string]Object)
		container:=Container{Name:accInfo.Container, Objs:objs}
		if len(account.Containers)==0{
			account.Containers=make(map[string]Container)
		}
		account.Containers[accInfo.Container]=container

		// marshall account
		accountBytes, err = account.MarshalMsg(nil)

		marshalledAcc:=MarshalledAcc{Bytes:accountBytes, Name:accInfo.AccName}

		err = ioutil.WriteFile(path+"/src/Account"+accInfo.AccName +strconv.Itoa(peerID),accountBytes,0777)

		// send marhshalled account to other nodes to save it
		var wg sync.WaitGroup
		wg.Add(len(accInfo.NodeList))
		var currentAddr int = rand.Intn(len(accInfo.NodeList))
		for _, peer := range accInfo.NodeList{
			peerURL := "http://" + peer[currentAddr] + "/SNPutAccP2PRequest"
			go func(p string, URL string){
				if peerID == int(p[len(p) - 1] - '0') {
					// Don't send to itself
				} else {
					r, w := io.Pipe()
					go func() {
						// save buffer to object
						err = json.NewEncoder(w).Encode(marshalledAcc)
						if err != nil {
							fmt.Println("Error encoding to pipe ", err.Error())
							return
						}
						defer w.Close()                        // close pipe when go routine finishes
					}()
					response, err := http.Post(peerURL, "application/json", r)
					if err != nil {
						fmt.Println("Error sending http POST p2p", err.Error())
					}
					responseCode:=response.StatusCode
					if responseCode == 201 {
						fmt.Println(response.StatusCode," SN created")
					}
					fmt.Println(responseCode)

				}
				wg.Done()
			}(peer[0], peerURL)
		}
		wg.Wait()
		httpVar.AccFileMutex.Unlock()



}


func checkAccCont(w http.ResponseWriter, r *http.Request){
	var accInfo AccInfo
	var peerID int = int(r.Host[len(r.Host) - 1] - '0')
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading ", err)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("error body ", err)
	}
	if err := json.Unmarshal(body, &accInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("error unmarshalling ", err)
		}
	}





	// read account from file
	accountBytes, err := ioutil.ReadFile(path+"/src/Account"+accInfo.AccName +strconv.Itoa(peerID))
	if err != nil {
		fmt.Println("SNPutCont Error: reading ", path+"/src/Account"+accInfo.AccName +strconv.Itoa(peerID))
	}
	account:=Account{}
	_,err = account.UnmarshalMsg(accountBytes)
	_, exists := account.Containers[accInfo.Container]
	if exists{
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusBadRequest)


}



func addObjToCont(w http.ResponseWriter, r *http.Request){
	var accInfo AccInfo
	var peerID int = int(r.Host[len(r.Host) - 1] - '0')
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading ", err)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("error body ", err)
	}
	if err := json.Unmarshal(body, &accInfo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		log.Println(err)
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("error unmarshalling ", err)
		}
	}





	// read account from file
	accountBytes, err := ioutil.ReadFile(path+"/src/Account"+accInfo.AccName +strconv.Itoa(peerID))
	if err != nil {
		fmt.Println("SNPutCont Error: reading ", path+"/src/Account"+accInfo.AccName +strconv.Itoa(peerID))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	account:=Account{}
	_,err = account.UnmarshalMsg(accountBytes)
	//add obj name to the container's object map
	if len(account.Containers[accInfo.Container].Objs)==0{
		fmt.Println("account.Containers[accInfo.Container].Objs len == 0")
		m:=make(map[string]Object)

		cont:=Container{Objs:m, Policy:account.Containers[accInfo.Container].Policy, Name:accInfo.Container}
		account.Containers[accInfo.Container]=cont
	}

	newObj:=Object{Name:accInfo.Obj, Size:accInfo.Size, PartsNum:accInfo.Parts, ParityNum:accInfo.Parity}
	account.Containers[accInfo.Container].Objs[accInfo.Obj]=newObj

	accountBytes, err = account.MarshalMsg(nil)
	if err != nil {
		fmt.Println("addObjToCont: error Marshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	marshalledAcc:=MarshalledAcc{Bytes:accountBytes, Name:accInfo.AccName}
	httpVar.AccFileMutex.Lock()

	err = ioutil.WriteFile(path+"/src/Account"+accInfo.AccName+strconv.Itoa(peerID),accountBytes,0777)
	var wg sync.WaitGroup
	wg.Add(len(accInfo.NodeList))
	var currentAddr int = rand.Intn(len(accInfo.NodeList))
	for _, peer := range accInfo.NodeList{
		peerURL := "http://" + peer[currentAddr] + "/SNPutAccP2PRequest"
		go func(p string, URL string){
			if peerID == int(p[len(p) - 1] - '0') {
				// Don't send to itself
			} else {
				r, w := io.Pipe()
				go func() {
					// save buffer to object
					err = json.NewEncoder(w).Encode(marshalledAcc)
					if err != nil {
						fmt.Println("Error encoding to pipe ", err.Error())
						return
					}
					defer w.Close()                        // close pipe //when go routine finishes
				}()
				response, err := http.Post(peerURL, "application/json", r)
				if err != nil {
					fmt.Println("Error sending http POST p2p", err.Error())
				}
				responseCode:=response.StatusCode
				if responseCode == 201 {
					fmt.Println(response.StatusCode," SN created")
				}
				fmt.Println(responseCode)

			}
			wg.Done()
		}(peer[0], peerURL)
	}
	wg.Wait()

	httpVar.AccFileMutex.Unlock()


	if err != nil {
		fmt.Println("addObjToCont: error Marshalling")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)




}




