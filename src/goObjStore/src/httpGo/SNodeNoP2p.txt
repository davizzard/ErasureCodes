package httpGo
import(
	"net/http"
	/*
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"encoding/json"
	"davizzard/ErasureCodes/src/goObjStore/src/httpVar"
	"strconv"
	"time"
*/
)
func SNodeListenNoP2P(w http.ResponseWriter, r *http.Request){
	/*
	var chunk msg
	// Get node ID
	var nodeID int =int(r.Host[len(r.Host)-1]-'0')

	// Listen to tracker
	if r.Method == http.MethodPost{

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("error reading ",err)
		}
		if err := r.Body.Close(); err != nil {
			fmt.Println("error body ",err)
		}
		if err := json.Unmarshal(body, &chunk); err != nil {
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(422) // unprocessable entity
			log.Println(err)
			if err := json.NewEncoder(w).Encode(err); err != nil {
				fmt.Println("error unmarshalling ",err)
			}
		}

		httpVar.DirMutex.Lock()
		// if data directory doesn't exist, create it
		_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/data")
		if err != nil {
			os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/data",0777)
		}

		// if data/chunk.Hash directory doesn't exist, create it
		_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/data/"+chunk.Hash)
		if err != nil {
			os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/data/"+chunk.Hash,0777)
		}

		// if data/chunk.Hash/nodeID directory doesn't exist, create it
		_, err = os.Stat(os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/data/"+chunk.Hash+"/"+strconv.Itoa( nodeID))
		if err != nil {
			err2:=os.Mkdir(os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/src/data/"+chunk.Hash+"/"+strconv.Itoa( nodeID),0777)
			if err2!=nil{
				fmt.Println("StorageNode error making dir", err.Error())
			}
		}
		httpVar.DirMutex.Unlock()

		// Save chunk to file
		err=ioutil.WriteFile(path+"/src/data/"+chunk.Hash+"/"+strconv.Itoa( nodeID)+"/NEW"+strconv.Itoa(httpVar.CurrentPart),[]byte(chunk.Text),0777)
		if err != nil {
			fmt.Println("StorageNodeListen: error creating/writing file", err.Error())
		}
		
		httpVar.TrackerMutex.Lock()
		httpVar.CurrentPart++
		httpVar.TrackerMutex.Unlock()


		if httpVar.CurrentPart == (totalPartsNum*chunk.Num)-1 {
			fmt.Println("..........................................Peer END ....................................................", time.Since(start))
		}
	}


*/
}
