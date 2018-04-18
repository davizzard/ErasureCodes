package httpGo
import(
	"net/http"
	"fmt"
	"os"
	"io"
	"github.com/davizzard/ErasureCodes/src/goObjStore/src/conf"
	"time"
	"sync"
	"strings"
	"encoding/json"
)

func PutObjAPI(w http.ResponseWriter, r *http.Request){
	var startPUT time.Time
	startPUT = time.Now()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		results:= strings.Split(r.URL.Path, "/")	// ["",account, container, object]
		addedResults:=results[1]+results[2]+results[3]

		// Creating temporary file to save received file (in the request body)
		file, err := os.Create(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/" + addedResults )
		CheckSimpleErr(err, nil, true)

		_, err = io.Copy(file, r.Body)
		CheckSimpleErr(err, nil, true)

		file.Close()

		// Creating a channel to control call
		putOK := make(chan bool)
		go PutObjProxy(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/" + addedResults, conf.TrackerAddr, conf.NumNodes, putOK ,results[1], results[2], results[3], addedResults)
		success := <-putOK
		if success == true {
			fmt.Println("put success ", time.Since(startPUT))
			w.WriteHeader(http.StatusCreated)
		} else {
			fmt.Println("put fail")
			w.WriteHeader(http.StatusBadRequest)
		}

		// Removing temporary file
		os.Remove(os.Getenv("GOPATH") + "/src/github.com/davizzard/ErasureCodes/src/goObjStore/src/" + addedResults )

	}()
	wg.Wait()
	fmt.Println("PUT: ",time.Since(startPUT))
}

func GetObjAPI(w http.ResponseWriter, r *http.Request){
	var startGET time.Time
	startGET = time.Now()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		results:= strings.Split(r.URL.Path, "/")	// ["",account, container, object]
		addedResults:=results[1]+results[2]+results[3]

		// Creating a channel to control call
		GetOK := make(chan bool)
		go GetObjProxy(addedResults, conf.ProxyAddr, conf.TrackerAddr, GetOK,results[1], results[2], results[3])
		success := <-GetOK
		if success == true {
			fmt.Println("get success ", time.Since(startGET))
			w.WriteHeader(http.StatusOK)
		} else {
			fmt.Println("get fail")
			w.WriteHeader(http.StatusBadRequest)
		}

	}()
	wg.Wait()
	fmt.Println("GET API: ",time.Since(startGET))
}


/*
func md5String(str string) string{
	hasher:=md5.New()
	_, err:= hasher.Write([]byte(str))
	if err != nil {
		fmt.Println(err)
	}
	return hex.EncodeToString(hasher.Sum(nil))
}
*/
func PutAccAPI(w http.ResponseWriter, r *http.Request){
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		accountName:=r.URL.Path[1:]
		if accountName==""{
			fmt.Println("create fail")
			w.WriteHeader(http.StatusBadRequest)
		} else{
			// Creating a channel to control call
			createOK := make(chan bool)
			go PutAccountProxy(accountName, createOK)
			success := <-createOK
			if success == true {
				fmt.Println("create success ")
				w.WriteHeader(http.StatusCreated)
			} else {
				fmt.Println("create fail")
				w.WriteHeader(http.StatusBadRequest)
			}
		}

	}()
	wg.Wait()
}


func GetAccAPI(w http.ResponseWriter, r *http.Request){
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		accountName:=r.URL.Path[1:]
		if accountName==""{
			fmt.Println("create fail")
			w.WriteHeader(http.StatusBadRequest)
		} else{
			account:= GetAccountProxy(accountName)

			// if account doesn't have Name means that couldn't find it
			if strings.Compare( account.Name, accountName)==0 {
				fmt.Println("get success ")

				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				// returning account
				if err := json.NewEncoder(w).Encode(account); err != nil {
					fmt.Println("GetNodes: error encoding response: ",err.Error())
				}
			} else {
				fmt.Println("get fail")
				w.WriteHeader(http.StatusBadRequest)
			}
		}

	}()
	wg.Wait()
}


func PutContAPI(w http.ResponseWriter, r *http.Request){
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		accountName:=r.URL.Path[1:]
		results:= strings.Split(accountName, "/")
		fmt.Println("PutContAPI: ",results[1])
		if accountName==""{
			fmt.Println("put fail")
			w.WriteHeader(http.StatusBadRequest)
		} else{
			// Creating a channel to control call
			createOK := make(chan bool)
			go PutContProxy(results[0], results[1], createOK)
			success := <-createOK
			if success == true {
				fmt.Println("put success ")
				w.WriteHeader(http.StatusCreated)
			} else {
				fmt.Println("put fail")
				w.WriteHeader(http.StatusBadRequest)
			}
		}

	}()
	wg.Wait()
}


func GetContAPI(w http.ResponseWriter, r *http.Request){
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		accountName:=r.URL.Path[1:]
		results:= strings.Split(accountName, "/")
		if accountName==""{
			fmt.Println("get fail")
			w.WriteHeader(http.StatusBadRequest)
		} else{

			container:=GetContProxy(results[0], results[1])

			// if container doesn't have Name means that couldn't find it
			if strings.Compare( container.Name,results[1] ) == 0 {
				fmt.Println("get success ")
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				err := json.NewEncoder(w).Encode(container)
				CheckJsonErr(err, nil, w)
			} else {
				fmt.Println("get fail",container.Name," /",results[1] )
				w.WriteHeader(http.StatusBadRequest)
			}
		}

	}()
	wg.Wait()
}



