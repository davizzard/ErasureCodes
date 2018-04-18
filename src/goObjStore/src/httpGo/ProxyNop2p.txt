package httpGo

import (
	/*
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
	"github.com/alruiz12/goObjStore/src/httpVar"
	"strings"
	"math/rand"
	*/

)
//var start time.Time
func PutNoP2P(filePath string, addr string, trackerAddr string, numNodes int){
	/*
	time.Sleep(1 * time.Second)
	start=time.Now()
	var hash string = md5sum(filePath)
	var err error

	// ask tracker for nodes
	quantityJson := `{"Quantity":"`+strconv.Itoa(numNodes)+`","Hash":"`+hash+`"}`
	reader := strings.NewReader(quantityJson)
	trackerURL:="http://"+trackerAddr+"/GetNodes"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("Put: error creating request: ",err.Error())
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Put: error sending request: ",err.Error())
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := res.Body.Close(); err != nil {
		panic(err)
	}
	var nodeList [][]string
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("Put: error unprocessable entity: ",err.Error())
		return
	}
	if err != nil {
		fmt.Println("Put: error reciving response: ",err.Error())
	}



	var currentPart int = 0
	var partSize int
	var currentNum int = 0
	var partBuffer []byte
	var writer *multipart.Writer
	var buf bytes.Buffer
	_,_=writer, buf // avoiding declared but not used

	var auxList []bool
	var i int = 0
	for i<numNodes {
		auxList=append(auxList, false)
		i++
	}
	httpVar.DirMutex.Lock()
	httpVar.HashMap[hash]=auxList
	httpVar.DirMutex.Unlock()
	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	defer file.Close()
	fileInfo, _ := file.Stat()
	text:=strconv.FormatInt(fileInfo.Size(),10)	// size
	size,_:=strconv.Atoi(text)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	totalPartsNum= int(math.Ceil(float64(size)/float64(fileChunk)))
//return
	for currentPart<totalPartsNum{
		partSize=int(math.Min(fileChunk, float64(size-(currentPart*fileChunk))))
		partBuffer=make([]byte,partSize)
		_,err = file.Read(partBuffer)		// Get chunk
		m:=msg{NodeList:nodeList, Num:numNodes, Hash:hash, Text:partBuffer, CurrentNode:currentNum, Name: currentPart}
 		m2:=msg{NodeList:nodeList, Num:numNodes, Hash:hash, Text:partBuffer, CurrentNode:currentNum, Name: currentPart}
 		m3:=msg{NodeList:nodeList, Num:numNodes, Hash:hash, Text:partBuffer, CurrentNode:currentNum, Name: currentPart}
		r, w :=io.Pipe()			// create pipe
		r2, w2 :=io.Pipe()                        // create pipe
		r3, w3 :=io.Pipe()                        // create pipe


			go func() {
				// save buffer to object
				err=json.NewEncoder(w).Encode(&m)
				if err != nil {
					fmt.Println("Error encoding to pipe ", err.Error())
				}
				defer  w.Close()                 // close pipe //when go routine finishes

			}()

		 go func() {
                                // save buffer to object
                                err=json.NewEncoder(w2).Encode(&m2)
                                if err != nil {
                                        fmt.Println("Error encoding to pipe ", err.Error())
                                }
                                defer  w2.Close()                 // close pipe //when go routine finishes

                        }()


		 go func() {
                                // save buffer to object
                                err=json.NewEncoder(w3).Encode(&m3)
                                if err != nil {
                                        fmt.Println("Error encoding to pipe ", err.Error())
                                }
                                defer  w3.Close()                 // close pipe //when go routine finishes

                        }()


		_, err := http.Post("http://"+nodeList[currentNum][0] + "/SNodeListenNoP2P", "application/json", r )
		if err != nil {
			fmt.Println("Error sending http POST ", err.Error())
		}
		currentNum=(currentNum+1)%numNodes


	         _, err = http.Post("http://"+nodeList[currentNum][0] + "/SNodeListenNoP2P", "application/json", r2 )
                 if err != nil {
                        fmt.Println("Error sending http POST ", err.Error())
                 }
                currentNum=(currentNum+1)%numNodes      
	
                
 		 _, err = http.Post("http://"+nodeList[currentNum][0] + "/SNodeListenNoP2P", "application/json", r3 )
                 if err != nil {
                        fmt.Println("Error sending http POST ", err.Error())
                 }


                currentPart++

		currentNum=(currentNum+1)%numNodes
	}
	fmt.Println("End of PutNoP2P!")
	*/
}



func GetNoP2P(Key string, proxyAddr []string, trackerAddr string){/*

	time.Sleep(1 * time.Second)

	// Ask tracker for nodes
	startGet=time.Now()
	var err error
	// ask tracker for nodes for a given key
	keyJson := `{"Key":"`+Key+`"}`
	reader := strings.NewReader(keyJson)
	trackerURL:="http://"+trackerAddr+"/GetNodesForKey"
	request, err := http.NewRequest("GET", trackerURL, reader)
	if err != nil {
		fmt.Println("Get: error creating request: ",err.Error())
	}
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Get: error sending request: ",err.Error())
	}
	body, err := ioutil.ReadAll(io.LimitReader(res.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := res.Body.Close(); err != nil {
		panic(err)
	}
	var nodeList []string
	if err := json.Unmarshal(body, &nodeList); err != nil {
		fmt.Println("Get: error unprocessable entity: ",err.Error())
		return
	}
	if err != nil {
		fmt.Println("Get: error reciving response: ",err.Error())
	}
	// Create folder for receiving
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/local",+0777)
	os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/local/"+Key,0777)

	node:=nodeList[rand.Intn(3)]
	// For each node ask for all their Proxy-pieces
	r, w :=io.Pipe()			// create pipe
	k:=jsonKeyURL{Key:Key, URL:proxyAddr[0]+"/ReturnData"}

	go func() {
		defer w.Close()			// close pipe when go routine finishes
		// save buffer to object
		err=json.NewEncoder(w).Encode(&k)
		if err != nil {
			fmt.Println("Error encoding to pipe ", err.Error())
		}
	}()
	url:="http://"+node+"/GetChunks"
	res, err = http.Post(url,"application/json", r )
	if err != nil {
		fmt.Println("Get2: error creating request: ",err.Error())
	}
	//fmt.Println("statusCode: ",res.StatusCode )
	if err := res.Body.Close(); err != nil {
		fmt.Println(err)
	}*/
}



