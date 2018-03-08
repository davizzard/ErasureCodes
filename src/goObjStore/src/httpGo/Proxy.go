package httpGo

import (
	"os"
	"fmt"
	"strconv"
	"math"
	"net/http"
	"io"
	"io/ioutil"
	"bytes"
	"mime/multipart"
	"time"
	"encoding/json"
	"crypto/md5"
	"encoding/hex"
	"github.com/davizzard/ErasureCodes/src/goObjStore/src/httpVar"
	"strings"
	"sync"
	"math/rand"
	"github.com/davizzard/ErasureCodes/src/goObjStore/src/conf"
	"github.com/davizzard/ErasureCodes/src/API"
)

const fileChunk = 1*(1<<10) // 1 KB
//const fileChunk = 4*(1<<20) // 4 MB
const parityShards = 2

type msg struct {
	NodeList [][]string
	Num int
	Hash string
	Text []byte 
	CurrentNode int
	Name int
	Size	int
	Parts	int
}
type putObjMSg struct {
	ID string
}
var startGet time.Time
func PutObjProxy(filePath string, trackerAddr string, numNodes int, putOK chan bool, account string, container string, objName string, fullName string) {
	time.Sleep(1 * time.Second)
	var err error
	var totalPartsNum int
	var size int

	// Get the account nodes
	var nodeList [][]string
	requestJson := `{"Quantity":"` + strconv.Itoa(conf.NumNodes) + `","ID":"` + account + `","Type":"account"}`
	reader := strings.NewReader(requestJson)
	trackerURL := "http://" + conf.TrackerAddr + "/GetNodesForKey"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("putObjProxy: error creating request: ", err.Error())
		putOK <- false
		return
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("putObjProxy: error sending request: ", err.Error())
		putOK <- false
		return
	}
	if res.StatusCode==http.StatusBadRequest{
		fmt.Println("PutObjProxy: ",res.StatusCode)
		putOK <- false
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		putOK <- false
		return
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
		putOK <- false
		return
	}
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("putObjProxy unmarshalling: error unprocessable entity: ", err.Error())
		putOK <- false
		return
	}



	// ask random node for correctness (account/container)
	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc := AccInfo{AccName:account, Container:container }

	r, w := io.Pipe()
	go func() {
		// save buffer to object
		err = json.NewEncoder(w).Encode(acc)
		if err != nil {
			fmt.Println("Error encoding to pipe ", err.Error())
			putOK <- false
			return
		}
		defer w.Close()                        // close pipe //when go routine finishes
	}()
	res, err = http.Post("http://" + nodeList[currentPeer][currentPeerAddr] + "/checkAccCont", "application/json", r)
	if err != nil {
		fmt.Println("Error sending http GET ", err.Error())
		putOK <- false
		return
	}
	if res.StatusCode == http.StatusBadRequest{
		fmt.Println("Error bad request")
		putOK <- false
		return
	}



	// Ask tracker for nodes to save object
	requestJson = `{"Quantity":"` + strconv.Itoa(numNodes) + `","ID":"` + fullName + `","Type":"object"}`
	reader = strings.NewReader(requestJson)
	trackerURL = "http://" + trackerAddr + "/GetNodes"
	request, err = http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("Put: error creating request: ", err.Error())
		putOK <- false
		return
	}
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Put: error sending request: ", err.Error())
		putOK <- false
		return
	}
	body, err = ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		putOK <- false
		return
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
		putOK <- false
		return
	}
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("Put: error unprocessable entity: ", err.Error())
		fmt.Println("NODELIST: ",nodeList)
		putOK <- false
		return
	}

	if len(nodeList)==0{
		fmt.Println(" no such name ")
		putOK <- false
		return
	}

	if err != nil {
		fmt.Println("Put: error receiving response: ", err.Error())
		putOK <- false
		return
	}
	var currentPart = 0
	var currentNum = 0
	var currentAdr = 0
	var fName string
	var writer *multipart.Writer
	var buf bytes.Buffer
	_, _ = writer, buf // avoiding declared but not used

	var auxList []bool
	var i int = 0
	for i < len(nodeList) {
		auxList = append(auxList, false)
		i++
	}

	/*
		// Open file containing the object
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Remove(filePath)
			file.Close()
			putOK <- false
			return
		}
		fileInfo, _ := file.Stat()
		text := strconv.FormatInt(fileInfo.Size(), 10)        // size
		size, _ := strconv.Atoi(text)
		if err != nil {
			fmt.Println(err.Error())
			os.Remove(filePath)
			file.Close()
			putOK <- false
			return
		}


		httpVar.TotalNumMutex.Lock()
		totalPartsNum = int(math.Ceil(float64(size) / float64(fileChunk)))
		httpVar.TotalNumMutex.Unlock()
		*/

	totalPartsNum, fName, size = API.EncodeFileAPI(filePath, fileChunk, parityShards, putOK)
	fmt.Println("EncodeFileApi Success!")

	if PrepareNodes(nodeList, fullName, currentNum, currentAdr, filePath) {
		putOK <- false
	}

	currentNum = 0

	// wait group will wait for all goroutines (one goroutine per chunk send)
	var wg sync.WaitGroup
	wg.Add(totalPartsNum)
	for currentPart < totalPartsNum {
		//partSize = int(math.Min(fileChunk, float64(size - (currentPart * fileChunk))))
		//partBuffer = make([]byte, partSize)
		//_, err = file.Read(partBuffer)                // Get chunk

		fileShardPath := conf.LocalDirectory + fName + "/" + strconv.Itoa(currentPart)
		//fileShardPath := "/root/Desktop/empty/go/src/davizzard/ErasureCodes/src/goObjStore/src/NEW" + strconv.Itoa(currentPart)

		if SendFileToNodes(fileShardPath, nodeList, fullName, currentNum, currentAdr, currentPart, &wg, filePath) {
			putOK <- false
		}

		currentPart++
		currentNum = (currentNum + 1) % len(nodeList)

		// Every 'numNodes' iterations, send chunk to next address, first send to different nodes, then change address
		if currentNum == 0 {
			currentAdr = (currentAdr + 1) % len(nodeList[currentNum])
		}
	}
	wg.Wait()

	// Add object to the container's object map
	currentPeer= rand.Intn(len(nodeList))
	currentPeerAddr = rand.Intn(len(nodeList))
	acc = AccInfo{AccName:account, Container:container, Obj:objName, Size:size, Parts:totalPartsNum, Parity:parityShards, NodeList:nodeList }
	r, w = io.Pipe()
	go func() {
		// save buffer to object
		err = json.NewEncoder(w).Encode(acc)
		if err != nil {
			fmt.Println("Error encoding to pipe ", err.Error())
			putOK <- false
			return
		}
		defer w.Close()                        // close pipe when go routine finishes
	}()
	resp, err := http.Post("http://" + nodeList[currentPeer][currentPeerAddr] + "/addObjToCont", "application/json", r)
	if err != nil {
		fmt.Println("Error sending http POST ", err.Error())
		putOK <- false
		return

	}
	if resp.StatusCode != http.StatusOK {
		putOK <- false
		return
	}

	err = os.RemoveAll(conf.LocalDirectory + "/" + fName)
	if err != nil {
		fmt.Println("error removing path and subDirs")
	}

	putOK <- true

	return
}



func PrepareNodes(nodeList [][]string, fullName string, currentNum int, address int, filePath string) (bool) {
	var exitStatus = false
	// Prepare nodes for content
	for currentNum < len(nodeList) {
		rpipe, wpipe := io.Pipe()
		mHash := putObjMSg{ID:fullName}
		go func() {
			// save buffer to object
			err := json.NewEncoder(wpipe).Encode(mHash)
			if err != nil {
				fmt.Println("Error encoding to pipe ", err.Error())
				if filePath != "" {
					os.Remove(filePath)
				}
				exitStatus = true
				return
			}
			defer wpipe.Close()                     // close pipe when go routine finishes
		}()

		_, err := http.Post("http://" + nodeList[currentNum][address] + "/prepSN", "application/json", rpipe)
		if err != nil {
			fmt.Println("to prepSN, Error sending http POST ", err.Error())
			if filePath != "" {
				os.Remove(filePath)
			}
			return true
		}
		currentNum++
	}
	return exitStatus
}


func SendFileToNodes(fileShardPath string, nodeList [][]string, fullName string, nodeNum int, address int, currentPart int, wg *sync.WaitGroup, filePath string) (bool) {
	var partSize int
	var partBuffer []byte
	var exitStatus = false

	fileShard, err := os.Open(fileShardPath)
	API.CheckErr(err)

	fileShardInfo, err := fileShard.Stat()
	API.CheckErr(err)

	text := strconv.FormatInt(fileShardInfo.Size(), 10) // size
	partSize, _ = strconv.Atoi(text)

	partBuffer = make([]byte, partSize)
	_, err = fileShard.Read(partBuffer)

	API.CheckErr(err)
	m := msg{NodeList:nodeList, Num:len(nodeList), Hash:fullName, Text:partBuffer, CurrentNode:nodeNum, Name: currentPart}
	go func(m2 msg, url string) {
		httpVar.SendReady <- 1
		r, w := io.Pipe()
		go func() {
			// save buffer to object
			err := json.NewEncoder(w).Encode(m2)
			if err != nil {
				fmt.Println("Error encoding to pipe ", err.Error())
				if filePath != "" {
					os.Remove(filePath)
				}
				fileShard.Close()
				exitStatus = true
				return
			}
			defer w.Close()                 // Close pipe //when go routine finishes
		}()

		_, err := http.Post(url, "application/json", r)
		if err != nil {
			fmt.Println("Error sending http POST ", err.Error())
			os.Remove(filePath)
			fileShard.Close()
			exitStatus = true
			return
		}
		defer wg.Done()
		<-httpVar.SendReady
	}(m, "http://" + nodeList[nodeNum][address] + "/SNPutObj")

	fileShard.Close()

	return exitStatus
	//wg.Wait()
		//os.Remove(fileShardPath)
	//wg.Wait()
	//os.Remove(filePath)
	//file.Close()

}


/*
md5sum opens the file we want to compute the hash and computes it
@param path to the file we want to split
returns the computed hash
*/
func md5sum(filePath string) string{
	file, err:=os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash,file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	mainFileHash:=hex.EncodeToString(hash.Sum(nil))
	return mainFileHash
}

type jsonKeyURL struct {
	Key string		`json:"Key"`
	URL string		`json:"URL"`
	Account string		`json:"Account"`
	Container string	`json:"Container"`
	Object string		`json:"Object"`
	GetID int		`json:"GetID"`
	NumParts int		`json:"NumParts"`
	NumParity int		`json:"NumParity"`
	NodeList [][]string `json: "NodeList"`
	ShardsInNodes int `json: "ShardsInNodes"`
}


func GetObjProxy(fullName string, proxyAddr []string, trackerAddr string, getOK chan bool, account string, container string, objName string,){
	fmt.Println("GetObjProxy init.")
	time.Sleep(1 * time.Second)
	startGet=time.Now()
	var err error
	// ask tracker for nodes for a given key
	requestJson := `{"ID":"`+ fullName +`","Type":"object"}`
	reader := strings.NewReader(requestJson)
	trackerURL:="http://"+trackerAddr+"/GetNodesForKey"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("GetObjProxy: error creating request: ",err.Error())
		getOK <- false
		return

	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("GetObjProxy: error sending request: ",err.Error())
		getOK <- false
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println("GetObjProxy ",err)
		getOK <- false
		return
		}
	if err := res.Body.Close(); err != nil {
		fmt.Println("GetObjProxy ",err)
		getOK <- false
		return
	}
	var nodeList [][]string
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("GetObjProxy ",err)
		getOK <- false
		return
	}

	// Create folder for receiving
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/local",+0777)
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/local/"+ fullName,0777)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	getID:=r.Int()
	httpVar.NumGetsMap[getID]=0

	acc := GetAccountProxy(account)
	numParts := acc.Containers[container].Objs[objName].PartsNum
	numParity := acc.Containers[container].Objs[objName].ParityNum

	var currentAddr int = rand.Intn(conf.PortsPerNode)

	httpVar.WaitingGroupNodes[getID] = new(sync.WaitGroup)

	// For each node ask for all their Proxy-pieces
	for index, node := range nodeList {
		r, w :=io.Pipe()			// create pipe
		k:=jsonKeyURL{Key:fullName, URL:proxyAddr[index]+"/ReturnObjProxy", Account: account, Container:container, Object:objName, GetID:getID, NumParts:numParts, NumParity:numParity, NodeList:nodeList}
		go func() {
			defer w.Close()			// close pipe when go routine finishes
			// save buffer to object
			err=json.NewEncoder(w).Encode(&k)
			if err != nil {
				fmt.Println("GetObjProxy: Error encoding to pipe ", err.Error())
				getOK <- false
				return

			}
		}()
		url:="http://"+node[currentAddr]+"/SNObjGetChunks"
		res, err := http.Post(url,"application/json", r )
		if err != nil {
			fmt.Println("GetObjProxy: error creating request: ",err.Error())
			getOK <- false
			return

		}
		if err := res.Body.Close(); err != nil {
			fmt.Println("GetObjProxy: ", err.Error())
			getOK <- false
			return
		}


	}

	// Waiting until all chunks are processed
	(*httpVar.WaitingGroupNodes[getID]).Wait()

	httpVar.GetMutex.Lock()
	numChunks := httpVar.NumGetsMap[getID]
	httpVar.GetMutex.Unlock()

	if numChunks >= numParts - numParity {
		httpVar.TotalNumMutex.Lock()
		fmt.Println("GatherPieces result:", GatherPieces(fullName, objName, numParts, numParity, nodeList))
		httpVar.TotalNumMutex.Unlock()
	} else {
		fmt.Printf("Unable to get object: %d shards found. Minimum %d shards needed.\n", numChunks, numParts - numParity)
	}

	httpVar.GetMutex.Lock()
	delete(httpVar.NumGetsMap, getID)
	httpVar.GetMutex.Unlock()

/*
	// Calling ReturnObjProxy to gather all the pieces
	_, w := io.Pipe()			// create pipe
	k:=jsonKeyURL{Key:fullName, URL:proxyAddr[0]+"/ReturnObjProxy", Account: account, Container:container, Object:objName, GetID:getID, NumParts:acc.Containers[container].Objs[objName].PartsNum, NumParity:acc.Containers[container].Objs[objName].ParityNum, NodeList:nodeList}
	go func() {
		defer w.Close()			// close pipe when go routine finishes
		// save buffer to object
		fmt.Println("HOLA!")
		err=json.NewEncoder(w).Encode(&k)
		if err != nil {
			fmt.Println("GetObjProxy: Error encoding to pipe ", err.Error())
			getOK <- false
			return

		}
	}()
*/
	fmt.Println("GetObjProxy end.")
	getOK <- true


}


func ReturnObjProxy(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ReturnObjProxy init.")
	var getmsg getMsg

	// Read request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error reading ", err)
	}
	if err := r.Body.Close(); err != nil {
		fmt.Println("error body ", err)
	}
	if err := json.Unmarshal(body, &getmsg); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		fmt.Println(err.Error())
		if err := json.NewEncoder(w).Encode(err); err != nil {
			fmt.Println("error unmarshalling ", err)
		}
	}

	// Update counter
	httpVar.GetMutex.Lock()
	httpVar.NumGetsMap[getmsg.GetID]++
	httpVar.GetMutex.Unlock()

	// Save chunk to file
	err = ioutil.WriteFile(path + "/src/local/" + getmsg.Key + "/" + getmsg.Name, getmsg.Text, 0777)
	if err != nil {
		fmt.Println("Peer: error creating/writing file p2p", err.Error())
	}

	fmt.Println("NumGetsMap: ", httpVar.NumGetsMap[getmsg.GetID])
	fmt.Println("getmsg.Parts: ", getmsg.Parts)

}


/*
CheckPieces walks through the subfiles directory, creates a new file to be filled out with the content of each subfile,
and compares the new hash with the original one.
@param path to the file we want to split
Returns true if both hash are identic and false if not
*/
func CheckPiecesObj(key string ,fileName string, filePath string, numNodes int, hash string) bool{
	 file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	text := strconv.FormatInt(fileInfo.Size(), 10)        // size
	size, _ := strconv.Atoi(text)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	totalPartsNumOriginal := int(math.Ceil(float64(size) / float64(fileChunk)))

	// Walking through subfiles directory
	path:=os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/data/"+key+"/"
	subDir, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return false
	}

	currentDir:=0
	for currentDir<numNodes{

		// Create new file to fill out
		_, err = os.Create(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src" + fileName+strconv.Itoa(currentDir))
		newFile, err := os.OpenFile(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src" + fileName+strconv.Itoa(currentDir), os.O_APPEND | os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println(err)
			return false
		}
		defer newFile.Close()

		files, err := ioutil.ReadDir(path+subDir[currentDir].Name() )

		// Trying to fill out the new file using subfiles (in order)
		var inOrderCount = 0
		var maxTimes int = 0
		var fileNameOriginal= fileName[:len(fileName)-4]

		for inOrderCount<totalPartsNumOriginal {
			for _, file := range files {
				if strings.Compare(file.Name(), fileNameOriginal + strconv.Itoa(inOrderCount)) == 0 || strings.Compare(file.Name(), "P2P" + strconv.Itoa(inOrderCount)) == 0{
					inOrderCount++
					//				fmt.Println(file.Name())
					currentFile, err := os.Open(path + subDir[currentDir].Name() +"/"+ file.Name())
					if err != nil {
						fmt.Println(err)
						return false
					}

					bytesCurrentFile, err := ioutil.ReadFile(path + subDir[currentDir].Name()+"/" +file.Name())

					_, err = newFile.WriteString(string(bytesCurrentFile))
					if err != nil {
						fmt.Println(err)
						return false
					}

					currentFile.Close()
				}

			}
			if inOrderCount == 0 {
				maxTimes++
			}
			if maxTimes > 1 {
				fmt.Println("maxTimes > 1 when looking for ", fileNameOriginal + strconv.Itoa(inOrderCount))
				return false
			}
		}

		// Compute and compare new hash
		newHash := md5sum(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src" + fileName+strconv.Itoa(currentDir))
		if strings.Compare(key, newHash) != 0 {
			return false
		}

		currentDir++
	}
	if currentDir==0{return false}	// Never got in loop

	// Checking Get output (locally)
	path=os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/local/"+key+"/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return false
	}

	// Create new file
	_, err = os.Create(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src" + fileName)
	newFile, err := os.OpenFile(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src" + fileName, os.O_APPEND | os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer newFile.Close()



	// Trying to fill out the new file using subfiles (in order)
	var inOrderCount = 0
	var maxTimes int = 0
	var fileNameOriginal= fileName[:len(fileName)-4]
	for inOrderCount<totalPartsNumOriginal {
		for _, file := range files {
			if strings.Compare(file.Name(), fileNameOriginal + strconv.Itoa(inOrderCount)) == 0 {
				inOrderCount++

				currentFile, err := os.Open(path + file.Name())
				if err != nil {
					fmt.Println(err)
					return false
				}

				bytesCurrentFile, err := ioutil.ReadFile(path + file.Name())

				_, err = newFile.WriteString(string(bytesCurrentFile))
				if err != nil {
					fmt.Println(err)
					return false
				}

				currentFile.Close()
			}

		}
		if inOrderCount == 0 {
			maxTimes++
		}
		if maxTimes > 1 {
			fmt.Println("maxTimes > 1 when looking for ", fileNameOriginal + strconv.Itoa(inOrderCount))
			return false
		}
	}

	// Compute and compare new hash
	newHash := md5sum(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src" + fileName)
	if strings.Compare(hash, newHash) != 0 {
		return false
	}

	return true
}

type AccInfo struct {
	NodeList    [][]string	`json:"NodeList"`
	Num         int		`json:"Num"`
	CurrentNode int		`json:"CurrentNode"`
	AccName     string	`json:"AccName"`
	Container   string	`json:"Container"`
	Obj         string	`json:"Obj"`
	Size	int		`json:"Size"`
	Parts	int		`json:"Parts"`
	Parity	int		`json:"Parity"`
}

func PutAccountProxy(name string, createOK chan bool){
	var nodeList [][]string
	var err error

	// Ask tracker for nodes
	requestJson := `{"Quantity":"` + strconv.Itoa(conf.NumNodes) + `","ID":"` + name + `","Type":"account"}`
	reader := strings.NewReader(requestJson)
	trackerURL := "http://" + conf.TrackerAddr + "/GetNodes"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("Put: error creating request: ", err.Error())
		createOK <- false
		return
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Put: error sending request: ", err.Error())
		createOK <- false
		return
	}
	if res.StatusCode==http.StatusBadRequest{
		fmt.Println(" no such name ")
		createOK <- false
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		createOK <- false
		return
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
		createOK <- false
		return
	}
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("Put: error unprocessable entity: ", err.Error())
		fmt.Println("CreateAccountProxy: NODELIST: ", nodeList)
		createOK <- false
		return
	}
	if err != nil {
		fmt.Println("Put: error receiving response: ", err.Error())
		createOK <- false
		return
	}
	if len(nodeList)==0{
		fmt.Println(" no such name ")
		createOK <- false
		return
	} else {
		fmt.Println("len ",len(nodeList))
	}

	// Randomly choose one Storage Node and one of its addresses to put account
	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc := AccInfo{NodeList:nodeList, Num:len(nodeList), CurrentNode:currentPeer, AccName:name }
	fmt.Print("JUST BEFORE SNPutAcc")
	r, w := io.Pipe()
	go func() {
		// save buffer to object
		err = json.NewEncoder(w).Encode(acc)
		if err != nil {
			fmt.Println("Error encoding to pipe ", err.Error())
			createOK <- false
			return
		}
		defer w.Close()                        // close pipe when go routine finishes
	}()
	_, err = http.Post("http://" + nodeList[currentPeer][currentPeerAddr] + "/SNPutAcc", "application/json", r)
	if err != nil {
		fmt.Println("Error sending http POST ", err.Error())
		createOK <- false
		return
	}
	createOK <- true

}

func GetAccountProxy(accountName string) Account{
	var nodeList [][]string
	var err error
	var account Account
	// Ask tracker for nodes
	requestJson := `{"Quantity":"` + strconv.Itoa(conf.NumNodes) + `","ID":"` + accountName + `","Type":"account"}`
	reader := strings.NewReader(requestJson)
	trackerURL := "http://" + conf.TrackerAddr + "/GetNodesForKey"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("putContProxy: error creating request: ", err.Error())
		return account
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("GetAccountProxy: error sending request: ", err.Error())
		return account
	}
	if res.StatusCode == http.StatusBadRequest{
		return account
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		return account
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
		return account
	}
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("GetAccountProxy: error unprocessable entity: ", err.Error())
		return account
	}

	// Randomly choose one Storage Node and one of its addresses to request account
	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc :=`{"Accname":"`+accountName+`"}`
	reader = strings.NewReader(acc)


	request, err = http.NewRequest("GET", "http://" + nodeList[currentPeer][currentPeerAddr] + "/SNGetAcc", reader)
	if err != nil {
		fmt.Println("Error sending http GET ", err.Error())
		return account
	}
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("GetAccountProxy: error sending request: ", err.Error())
		return account
	}
	if res.StatusCode == http.StatusBadRequest{
		return account
	}
	body, err = ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		return account
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
		return account
	}

	if err := json.Unmarshal(body, &account); err != nil {
		fmt.Println("GetAccountProxy: error unprocessable entity: ", err.Error())
		return account
	}
	return account
}

func PutContProxy(account string, container string, createOK chan bool){
	var nodeList [][]string
	var err error

	// Ask tracker for nodes
	requestJson := `{"Quantity":"` + strconv.Itoa(conf.NumNodes) + `","ID":"` + account + `","Type":"account"}`
	reader := strings.NewReader(requestJson)
	trackerURL := "http://" + conf.TrackerAddr + "/GetNodesForKey"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("putContProxy: error creating request: ", err.Error())
		createOK <- false
		return
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("putContProxy: error sending request: ", err.Error())
		createOK <- false
		return
	}
	if res.StatusCode == http.StatusBadRequest{
		createOK <- false
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		createOK <- false
		return
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
		createOK <- false
		return
	}
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("putContProxy: error unprocessable entity: ", err.Error())
		createOK <- false
		return
	}

	// Randomly choose one Storage Node and one of its addresses to put container
	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc := AccInfo{NodeList:nodeList, Num:len(nodeList), CurrentNode:currentPeer, AccName:account, Container:container }

	r, w := io.Pipe()
	go func() {
		// save buffer to object
		err = json.NewEncoder(w).Encode(acc)
		if err != nil {
			fmt.Println("Error encoding to pipe ", err.Error())
			createOK <- false
			return
		}
		defer w.Close()                        // close pipe //when go routine finishes
	}()
	_, err = http.Post("http://" + nodeList[currentPeer][currentPeerAddr] + "/SNPutCont", "application/json", r)
	if err != nil {
		fmt.Println("Error sending http POST ", err.Error())
		createOK <- false
		return
	}

	createOK <- true
}



func GetContProxy(accountName string, containerName string) Container{
	var nodeList [][]string
	var err error
	var account Account
	var container Container
	// Ask tracker for nodes
	requestJson := `{"Quantity":"` + strconv.Itoa(conf.NumNodes) + `","ID":"` + accountName + `","Type":"account"}`
	reader := strings.NewReader(requestJson)
	trackerURL := "http://" + conf.TrackerAddr + "/GetNodesForKey"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("putContProxy: error creating request: ", err.Error())
		return container
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("GetAccountProxy: error sending request: ", err.Error())
		return container
	}
	if res.StatusCode == http.StatusBadRequest{
		return container
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		return container
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
		return container
	}
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("GetAccountProxy: error unprocessable entity: ", err.Error())
		return container
	}

	// Randomly choose one Storage Node and one of its addresses to request container
	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc :=`{"Accname":"`+accountName+`"}`
	reader = strings.NewReader(acc)


	request, err = http.NewRequest("GET", "http://" + nodeList[currentPeer][currentPeerAddr] + "/SNGetAcc", reader)
	if err != nil {
		fmt.Println("Error sending http GET ", err.Error())
		return container
	}
	res, err = http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("GetAccountProxy: error sending request: ", err.Error())
		return container
	}
	if res.StatusCode == http.StatusBadRequest{
		return container
	}
	body, err = ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		fmt.Println(err)
		return container
	}
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
		return container
	}

	if err := json.Unmarshal(body, &account); err != nil {
		fmt.Println("GetAccountProxy: error unprocessable entity: ", err.Error())
		return container
	}

	return account.Containers[containerName]
}



func CheckFileReplication(fileType string, name string, replication int) bool{
	if replication<2 {return false}
	var path = (os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/")
	 hashList := make([]string, replication)
	var currentReplica int = 0
	var i int =0
	for currentReplica < replication && i<10{
		_, err := os.Stat(path+fileType+name+strconv.Itoa(currentReplica+1))
		if err != nil {
			fmt.Println("Error: "+path+fileType+name+strconv.Itoa(currentReplica+1)+" does not exist")
		} else{
			hashList[currentReplica] = md5sum(path+fileType+name+strconv.Itoa(currentReplica+1) )
			if strings.Compare( hashList[currentReplica], string(""))==0{
				return false
			}
			currentReplica++
		}

		i++
	}
	if currentReplica!=replication{
		fmt.Println("replica missing")
		return false
	}
	firstHash:=hashList[0]
	currentReplica = 1
	for currentReplica < replication{
		if strings.Compare( hashList[currentReplica], firstHash) != 0{
			return false
		}
		currentReplica++
	}

	return true
}


/*
CheckPieces walks through the subfiles directory, creates a new file to be filled out with the content of each subfile,
and compares the new hash with the original one.
@param path to the file we want to split
Returns true if both hash are identic and false if not
*/

func GatherPieces(key string , objName string, totalParts int, parityShards int, nodeList [][]string) bool{

	/*
	// Checking Get output (locally)
	path=os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/local/"+key+"/"
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println("Readir",err, " ",path)
		return false
	}

	// Create new file
	_, err = os.Create(conf.DownloadsDirectory + key)
	if err != nil {
		fmt.Println(" create ",err)
	}
	newFile, err := os.OpenFile(conf.DownloadsDirectory + key, os.O_APPEND | os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("GatherPieces",err)
		return false
	}
	defer newFile.Close()


	// Trying to fill out the new file using subfiles (in order)
	var inOrderCount = 0
	var maxTimes int = 0

	var fileNameOriginal= conf.LocalFileName[:len(conf.LocalFileName)-4]
	for inOrderCount<totalParts {
		for _, file := range files {
			if strings.Compare(file.Name(), fileNameOriginal + strconv.Itoa(inOrderCount)) == 0 {
				inOrderCount++
				currentFile, err := os.Open(path + file.Name())
				if err != nil {
					fmt.Println(err)
					return false
				}

				bytesCurrentFile, err := ioutil.ReadFile(path + file.Name())

				_, err = newFile.WriteString(string(bytesCurrentFile))
				if err != nil {
					fmt.Println(err)
					return false
				}

				currentFile.Close()
			}

		}
		if inOrderCount == 0 {
			maxTimes++
		}
		if maxTimes > 1 {
			fmt.Println("maxTimes > 1 when looking for ", fileNameOriginal + strconv.Itoa(inOrderCount))
			return false
		}

	}
	*/
	fmt.Println("GATHER PIECES")
	// Checking Get output (locally)
	path=os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/local/"+key+"/"

	//os.Remove(path+"NEW0")
	//os.Remove(path+"NEW18")
	//fmt.Println("IN PROXY: FILES REMOVED.")

	API.DecodeFileAPI(path, key, totalParts-parityShards, parityShards, conf.ChunkProxyName, objName, nodeList)

	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println("error removing path and subDirs")
	}
	return true
}











