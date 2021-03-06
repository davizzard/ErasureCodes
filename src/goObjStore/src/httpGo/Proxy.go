package httpGo

import (
	"os"
	"fmt"
	"strconv"
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
	CheckSimpleErr(err, putOK, true)

	res, err := http.DefaultClient.Do(request)
	CheckSimpleErr(err, putOK, true)

	if res.StatusCode==http.StatusBadRequest{
		fmt.Println("PutObjProxy: ",res.StatusCode)
		putOK <- false
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	CheckSimpleErr(err, putOK, true)

	err = res.Body.Close()
	CheckSimpleErr(err, putOK, true)

	err = json.Unmarshal(body, &nodeList)
	CheckSimpleErr(err, putOK, true)



	// ask random node for correctness (account/container)
	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc := AccInfo{AccName:account, Container:container }

	r, w := io.Pipe()
	go func() {
		// save buffer to object
		err = json.NewEncoder(w).Encode(acc)
		CheckSimpleErr(err, putOK, true)
		defer w.Close()                        // close pipe //when go routine finishes
	}()
	res, err = http.Post("http://" + nodeList[currentPeer][currentPeerAddr] + "/checkAccCont", "application/json", r)
	CheckSimpleErr(err, putOK, true)
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
	CheckSimpleErr(err, putOK, true)
	res, err = http.DefaultClient.Do(request)

	CheckSimpleErr(err, putOK, true)
	body, err = ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	CheckSimpleErr(err, putOK, true)

	err = res.Body.Close()
	CheckSimpleErr(err, putOK, true)

	err = json.Unmarshal(body, &nodeList)
	CheckSimpleErr(err, putOK, true)

	CheckLengthErr(len(nodeList), " no such name ", putOK, true)

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

	totalPartsNum, fName, size = EncodeFileAPI(filePath, fileChunk, parityShards, putOK)
	fmt.Println("EncodeFileApi Success!")

	PrepareNodes(nodeList, fullName, currentNum, currentAdr, filePath)

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

		SendFileToNodes(fileShardPath, nodeList, fullName, currentNum, currentAdr, currentPart, &wg, filePath)

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
		CheckSimpleErr(err, putOK, true)
		defer w.Close()                        // close pipe when go routine finishes
	}()
	resp, err := http.Post("http://" + nodeList[currentPeer][currentPeerAddr] + "/addObjToCont", "application/json", r)
	CheckSimpleErr(err, putOK, true)

	if resp.StatusCode != http.StatusOK {
		putOK <- false
		return
	}

	err = os.RemoveAll(conf.LocalDirectory + "/" + fName)
	CheckSimpleErr(err, nil, false)

	putOK <- true

	return
}



func PrepareNodes(nodeList [][]string, fullName string, currentNum int, address int, filePath string) {
	// Prepare nodes for content
	for currentNum < len(nodeList) {
		rpipe, wpipe := io.Pipe()
		mHash := putObjMSg{ID:fullName}
		go func() {
			// save buffer to object
			err := json.NewEncoder(wpipe).Encode(mHash)
			CheckSimpleErr(err, nil, true)
			defer wpipe.Close()                     // close pipe when go routine finishes
		}()

		_, err := http.Post("http://" + nodeList[currentNum][address] + "/prepSN", "application/json", rpipe)
		CheckSimpleErr(err, nil, true)
		currentNum++
	}
}


func SendFileToNodes(fileShardPath string, nodeList [][]string, fullName string, nodeNum int, address int, currentPart int, wg *sync.WaitGroup, filePath string) {
	var partSize int
	var partBuffer []byte

	fileShard, err := os.Open(fileShardPath)
	CheckSimpleErr(err, nil, true)

	fileShardInfo, err := fileShard.Stat()
	CheckSimpleErr(err, nil, true)

	text := strconv.FormatInt(fileShardInfo.Size(), 10) // size
	partSize, _ = strconv.Atoi(text)

	partBuffer = make([]byte, partSize)
	_, err = fileShard.Read(partBuffer)

	CheckSimpleErr(err, nil, true)
	m := msg{NodeList:nodeList, Num:len(nodeList), Hash:fullName, Text:partBuffer, CurrentNode:nodeNum, Name: currentPart}
	go func(m2 msg, url string) {
		httpVar.SendReady <- 1
		r, w := io.Pipe()
		go func() {
			// save buffer to object
			err := json.NewEncoder(w).Encode(m2)
			CheckComplexErr(err, nil, filePath, *fileShard, true)
			defer w.Close()                 // Close pipe //when go routine finishes
		}()

		_, err := http.Post(url, "application/json", r)
		CheckComplexErr(err, nil, filePath, *fileShard, true)
		defer wg.Done()
		<-httpVar.SendReady
	}(m, "http://" + nodeList[nodeNum][address] + "/SNPutObj")

	fileShard.Close()

}


/*
md5sum opens the file we want to compute the hash and computes it
@param path to the file we want to split
returns the computed hash
*/
func md5sum(filePath string) string{
	file, err:=os.Open(filePath)
	CheckSimpleErr(err, nil, false)
	defer file.Close()

	hash := md5.New()
	_, err = io.Copy(hash,file)
	CheckSimpleErr(err, nil, false)
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
}


func GetObjProxy(fullName string, proxyAddr []string, trackerAddr string, getOK chan bool, account string, container string, objName string,){
	time.Sleep(1 * time.Second)
	startGet=time.Now()
	var err error
	// ask tracker for nodes for a given key
	requestJson := `{"ID":"`+ fullName +`","Type":"object"}`
	reader := strings.NewReader(requestJson)
	trackerURL:="http://"+trackerAddr+"/GetNodesForKey"
	request, err := http.NewRequest("GET", trackerURL, reader)
	CheckSimpleErr(err, getOK, true)

	res, err := http.DefaultClient.Do(request)
	CheckSimpleErr(err, getOK, true)

	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	CheckSimpleErr(err, getOK, true)

	err = res.Body.Close()
	CheckSimpleErr(err, getOK, true)

	var nodeList [][]string
	err = json.Unmarshal(body, &nodeList)
	CheckSimpleErr(err, getOK, true)

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
			CheckSimpleErr(err, getOK, true)
		}()
		url:="http://"+node[currentAddr]+"/SNObjGetChunks"
		res, err := http.Post(url,"application/json", r )
		CheckSimpleErr(err, getOK, true)

		err = res.Body.Close()
		CheckSimpleErr(err, getOK, true)

	}

	// Waiting until all chunks are processed
	(*httpVar.WaitingGroupNodes[getID]).Wait()

	httpVar.GetMutex.Lock()
	numChunks := httpVar.NumGetsMap[getID]
	httpVar.GetMutex.Unlock()

	if numChunks >= numParts - numParity {
		httpVar.TotalNumMutex.Lock()
		fmt.Println("GatherPieces result:", GatherPieces(fullName, numParts, numParity, nodeList))
		httpVar.TotalNumMutex.Unlock()
	} else {
		fmt.Printf("Unable to get object: %d shards found. Minimum %d shards needed.\n", numChunks, numParts - numParity)
		getOK <- false
	}

	httpVar.GetMutex.Lock()
	delete(httpVar.NumGetsMap, getID)
	httpVar.GetMutex.Unlock()

	fmt.Println("GetObjProxy end.")
	getOK <- true


}


func ReturnObjProxy(w http.ResponseWriter, r *http.Request) {
	var getmsg getMsg

	// Read request
	body, err := ioutil.ReadAll(r.Body)
	CheckSimpleErr(err, nil, true)

	err = r.Body.Close()
	CheckSimpleErr(err, nil, true)

	err = json.Unmarshal(body, &getmsg)
	CheckJsonErr(err, nil, w)

	// Update counter
	httpVar.GetMutex.Lock()
	httpVar.NumGetsMap[getmsg.GetID]++
	httpVar.GetMutex.Unlock()

	// Save chunk to file
	err = ioutil.WriteFile(path + "/src/local/" + getmsg.Key + "/" + getmsg.Name, getmsg.Text, 0777)
	CheckSimpleErr(err, nil, true)

}


/*
CheckPieces walks through the subfiles directory, creates a new file to be filled out with the content of each subfile,
and compares the new hash with the original one.
@param path to the file we want to split
Returns true if both hash are identic and false if not

func CheckPiecesObj(key string ,fileName string, filePath string, numNodes int, hash string) bool{
	 file, err := os.Open(filePath)
	CheckSimpleErr(err, nil, true)
	defer file.Close()

	fileInfo, _ := file.Stat()
	text := strconv.FormatInt(fileInfo.Size(), 10)        // size
	size, _ := strconv.Atoi(text)
	CheckSimpleErr(err, nil, true)

	totalPartsNumOriginal := int(math.Ceil(float64(size) / float64(fileChunk)))

	// Walking through subfiles directory
	path:=os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/data/"+key+"/"
	subDir, err := ioutil.ReadDir(path)
	CheckSimpleErr(err, nil, true)

	currentDir:=0
	for currentDir<numNodes{

		// Create new file to fill out
		_, err = os.Create(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src" + fileName+strconv.Itoa(currentDir))
		newFile, err := os.OpenFile(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src" + fileName+strconv.Itoa(currentDir), os.O_APPEND | os.O_WRONLY, 0666)
		CheckSimpleErr(err, nil, true)
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
					CheckSimpleErr(err, nil, true)

					bytesCurrentFile, err := ioutil.ReadFile(path + subDir[currentDir].Name()+"/" +file.Name())

					_, err = newFile.WriteString(string(bytesCurrentFile))
					CheckSimpleErr(err, nil, true)

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
	CheckSimpleErr(err, nil, true)

	// Create new file
	_, err = os.Create(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src" + fileName)
	newFile, err := os.OpenFile(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src" + fileName, os.O_APPEND | os.O_WRONLY, 0666)
	CheckSimpleErr(err, nil, true)
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
				CheckSimpleErr(err, nil, true)

				bytesCurrentFile, err := ioutil.ReadFile(path + file.Name())

				_, err = newFile.WriteString(string(bytesCurrentFile))
				CheckSimpleErr(err, nil, true)

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
*/

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
	CheckSimpleErr(err, createOK, true)

	res, err := http.DefaultClient.Do(request)
	CheckSimpleErr(err, createOK, true)

	if res.StatusCode==http.StatusBadRequest{
		fmt.Println(" no such name ")
		createOK <- false
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	CheckSimpleErr(err, createOK, true)

	err = res.Body.Close()
	CheckSimpleErr(err, createOK, true)

	err = json.Unmarshal(body, &nodeList)
	CheckSimpleErr(err, createOK, true)

	CheckLengthErr(len(nodeList), " no such name ", createOK, true)

	// Randomly choose one Storage Node and one of its addresses to put account
	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc := AccInfo{NodeList:nodeList, Num:len(nodeList), CurrentNode:currentPeer, AccName:name }
	fmt.Print("JUST BEFORE SNPutAcc")
	r, w := io.Pipe()
	go func() {
		// save buffer to object
		err = json.NewEncoder(w).Encode(acc)
		CheckSimpleErr(err, createOK, true)

		defer w.Close()                        // close pipe when go routine finishes
	}()
	_, err = http.Post("http://" + nodeList[currentPeer][currentPeerAddr] + "/SNPutAcc", "application/json", r)
	CheckSimpleErr(err, createOK, true)

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
	CheckSimpleErr(err, nil, true)

	res, err := http.DefaultClient.Do(request)
	CheckSimpleErr(err, nil, true)

	if res.StatusCode == http.StatusBadRequest{
		return account
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	CheckSimpleErr(err, nil, true)

	err = res.Body.Close()
	CheckSimpleErr(err, nil, true)

	err = json.Unmarshal(body, &nodeList)
	CheckSimpleErr(err, nil, true)

	// Randomly choose one Storage Node and one of its addresses to request account
	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc :=`{"Accname":"`+accountName+`"}`
	reader = strings.NewReader(acc)


	request, err = http.NewRequest("GET", "http://" + nodeList[currentPeer][currentPeerAddr] + "/SNGetAcc", reader)
	CheckSimpleErr(err, nil, true)

	res, err = http.DefaultClient.Do(request)
	CheckSimpleErr(err, nil, true)

	if res.StatusCode == http.StatusBadRequest{
		return account
	}
	body, err = ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	CheckSimpleErr(err, nil, true)

	err = res.Body.Close()
	CheckSimpleErr(err, nil, true)

	err = json.Unmarshal(body, &account)
	CheckSimpleErr(err, nil, true)

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
	CheckSimpleErr(err, createOK, true)

	res, err := http.DefaultClient.Do(request)
	CheckSimpleErr(err, createOK, true)

	if res.StatusCode == http.StatusBadRequest{
		createOK <- false
		return
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	CheckSimpleErr(err, createOK, true)

	err = res.Body.Close()
	CheckSimpleErr(err, createOK, true)

	err = json.Unmarshal(body, &nodeList)
	CheckSimpleErr(err, nil, true)

	// Randomly choose one Storage Node and one of its addresses to put container
	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc := AccInfo{NodeList:nodeList, Num:len(nodeList), CurrentNode:currentPeer, AccName:account, Container:container }

	r, w := io.Pipe()
	go func() {
		// save buffer to object
		err = json.NewEncoder(w).Encode(acc)
		CheckSimpleErr(err, createOK, true)
		defer w.Close()                        // close pipe //when go routine finishes
	}()
	_, err = http.Post("http://" + nodeList[currentPeer][currentPeerAddr] + "/SNPutCont", "application/json", r)
	CheckSimpleErr(err, createOK, true)

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
	CheckSimpleErr(err, nil, true)

	res, err := http.DefaultClient.Do(request)
	CheckSimpleErr(err, nil, true)

	if res.StatusCode == http.StatusBadRequest{
		return container
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	CheckSimpleErr(err, nil, true)

	err = res.Body.Close()
	CheckSimpleErr(err, nil, true)

	err = json.Unmarshal(body, &nodeList)
	CheckSimpleErr(err, nil, true)

	// Randomly choose one Storage Node and one of its addresses to request container
	currentPeer:= rand.Intn(len(nodeList))
	currentPeerAddr := rand.Intn(len(nodeList))
	acc :=`{"Accname":"`+accountName+`"}`
	reader = strings.NewReader(acc)


	request, err = http.NewRequest("GET", "http://" + nodeList[currentPeer][currentPeerAddr] + "/SNGetAcc", reader)
	CheckSimpleErr(err, nil, true)

	res, err = http.DefaultClient.Do(request)
	CheckSimpleErr(err, nil, true)

	if res.StatusCode == http.StatusBadRequest{
		return container
	}
	body, err = ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	CheckSimpleErr(err, nil, true)

	err = res.Body.Close()
	CheckSimpleErr(err, nil, true)

	err = json.Unmarshal(body, &account)
	CheckSimpleErr(err, nil, true)

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


func GatherPieces(key string , totalParts int, parityShards int, nodeList [][]string) bool {
	var exitStatus = false
	fmt.Println("Gather Pieces init.")
	path=os.Getenv("GOPATH")+"/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/local/"+key+"/"

	if DecodeFileAPI(path, key, totalParts-parityShards, parityShards, conf.ChunkProxyName, nodeList) {
		exitStatus = true
	}

	err := os.RemoveAll(path)
	CheckSimpleErr(err, nil, false)

	return exitStatus
}











