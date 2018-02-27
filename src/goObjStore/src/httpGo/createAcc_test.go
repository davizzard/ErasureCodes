package httpGo
import (
	"testing"
	"fmt"
	"net/http"
	"time"
	"os"
	"os/exec"
	"strings"
	"davizzard/ErasureCodes/src/goObjStore/src/conf"
)
func TestCreateAccountAPI(t *testing.T) {
	router :=MyNewRouter()
	const IP = "127.0.0.1"
	var path = os.Getenv("GOPATH")+"/src/davizzard/ErasureCodes/src/goObjStore/src"

	//var filePath = os.Getenv("GOPATH")+"/src/github.com/alruiz12/goObjStore/dataset.xml"

	// TRACKER
	var trackerAddr = IP+":8000"

	// PEERS
	// IMPORTANT NOTE: last character of port is the peer's internal identifier
	var Peer1a = IP+":8011"
	var Peer1b = IP+":8021"
	var Peer1c = IP+":8031"
	var peer1List = []string{Peer1a, Peer1b, Peer1c}

	var Peer2a = IP+":8012"
	var Peer2b = IP+":8022"
	var Peer2c = IP+":8032"
	var peer2List = []string{Peer2a, Peer2b, Peer2c}

	var Peer3a = IP+":8013"
	var Peer3b = IP+":8023"
	var Peer3c = IP+":8033"
	var peer3List = []string{Peer3a, Peer3b, Peer3c}

	var Peers =[][]string{peer1List, peer2List, peer3List/*, peer4List, peer5List*/}

	// PROXY
	var proxy1 = IP+":8070"
	var proxy2 = IP+":8071"
	var proxy3 = IP+":8072"


	peer1arun:=&http.Server{Addr:Peer1a, Handler:router}
	peer1brun:=&http.Server{Addr:Peer1b, Handler:router}
	peer1crun:=&http.Server{Addr:Peer1c, Handler:router}
	peer2arun:=&http.Server{Addr:Peer2a, Handler:router}
	peer2brun:=&http.Server{Addr:Peer2b, Handler:router}
	peer2crun:=&http.Server{Addr:Peer2c, Handler:router}
	peer3arun:=&http.Server{Addr:Peer3a, Handler:router}
	peer3brun:=&http.Server{Addr:Peer3b, Handler:router}
	peer3crun:=&http.Server{Addr:Peer3c, Handler:router}
	/*peer4:=&http.Server{Addr:"8084", Handler:router}
	peer5:=&http.Server{Addr:"8085", Handler:router}
	*/

	proxy1run:=&http.Server{Addr:proxy1, Handler:router}
	proxy2run:=&http.Server{Addr:proxy2, Handler:router}
	proxy3run:=&http.Server{Addr:proxy3, Handler:router}
	/*proxy4:=&http.Server{Addr:"8074", Handler:router}
	proxy5:=&http.Server{Addr:"8075", Handler:router}
	*/



	tracker:=&http.Server{Addr:trackerAddr, Handler:router}

	go func(){
		if err := peer1arun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer1brun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer1crun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()

	go func(){
		if err := peer2arun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer2brun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer2crun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()

	go func(){
		if err := peer3arun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer3brun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := peer3crun.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()

	go func(){
		if err := proxy1run.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := proxy2run.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := proxy3run.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()
	go func(){
		if err := tracker.ListenAndServe(); err!=nil{
			fmt.Println("ListenAndServe error", err.Error())
		}
	}()


	StartTracker(Peers)

	var (
		cmdOut []byte
		err    error
	)
	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlCreateAccSuccess.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp := string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"201")==0 {
		fmt.Println(resp+ " created")
	}else{t.Error("Account not created")}
	fmt.Println("---------------------------------------------------")


	req, err := http.NewRequest("PUT", "http://localhost:8000/alvaroGO", nil)
	if err != nil {
		t.Error(" error creating post request to http://localhost:8000/alvaroGO")
	}
	response , err := http.DefaultClient.Do(req)
	if err != nil {
		t.Error("Error doing request")
	}
	fmt.Println("response: ",response.StatusCode )
	if response.StatusCode == 201 {
		fmt.Println(response.StatusCode," created")
	} else{
		t.Error("Account not created")
	}
	fmt.Println("---------------------------------------------------")


	fileReplicationCorrect:=CheckFileReplication("Account","alvaro",3)
	if fileReplicationCorrect==false{
		t.Error("file replication not correct")
	}
	fmt.Println(fileReplicationCorrect)
	fmt.Println("---------------------------------------------------")


	fileReplicationCorrect=CheckFileReplication("Account","alvaroGO",3)
	if fileReplicationCorrect==false{
		t.Error("file replication not correct")
	}
	fmt.Println(fileReplicationCorrect)
	fmt.Println("---------------------------------------------------")


	if CheckFileReplication("Account","accountWillFail", 1)==true{
		t.Error("file replication not correct")
	} else {fmt.Println("false")}
	fmt.Println("---------------------------------------------------")


	if CheckFileReplication("Account", "alvaro", 4)==true{
		t.Error("file replication not correct")
	} else {fmt.Println("false")}
	fmt.Println("---------------------------------------------------")


	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlPutContSuccess.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"201")==0 {
		fmt.Println(resp+ " created")
	}else{t.Error("Obj not created")}
	fmt.Println("---------------------------------------------------")


	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlPutContFailure.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"400")==0 {
		fmt.Println(resp+ " failed as expected")
	}else{t.Error("Expecting failure")}
	fmt.Println("---------------------------------------------------")



	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlPutObjSuccess.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"201")==0 {
		fmt.Println(resp+ " created")
	}else{t.Error("Obj not created")}
	fmt.Println("---------------------------------------------------")


	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlPutObjFailure.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"400")==0 {
		fmt.Println(resp+ " failed as expected")
	}else{t.Error("Expecting failure")}
	fmt.Println("---------------------------------------------------")

	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlPutObjFailure2.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"400")==0 {
		fmt.Println(resp+ " failed as expected")
	}else{t.Error("Expecting failure")}
	fmt.Println("---------------------------------------------------")

	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlPutObjFailure3.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"400")==0 {
		fmt.Println(resp+ " failed as expected")
	}else{t.Error("Expecting failure")}
	fmt.Println("---------------------------------------------------")



	req, err = http.NewRequest("PUT", "http://localhost:8000/InventedAccount/contB2B", nil)
	if err != nil {
		t.Error(" error creating post request to http://localhost:8000/ACCgo3/cont1")
	}
	response , err = http.DefaultClient.Do(req)
	if err != nil {
		t.Error("Error doing request")
	}
	fmt.Println("response: ",response.StatusCode )
	if response.StatusCode == 400 {
		fmt.Println(response.StatusCode," bad request as expected")
	} else{
		t.Error("expecting failure")
	}
	fmt.Println("---------------------------------------------------")


	// Same Obj Again
	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlPutObjSuccess.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"201")==0 {
		fmt.Println(resp+ " created")
	}else{t.Error("Obj not created")}
	fmt.Println("---------------------------------------------------")



	fileReplicationCorrect=CheckFileReplication("Account","alvaro",3)
	if fileReplicationCorrect==false{
		t.Error("file replication not correct")
	}
	fmt.Println(fileReplicationCorrect)
	fmt.Println("---------------------------------------------------")


	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlGetObjSuccess.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"200")==0 {
		fmt.Println(resp+ " OK")
	}else{t.Error("Obj not obtained")}
	fmt.Println("---------------------------------------------------")

	time.Sleep(2* time.Second)

	currentHash := md5sum(conf.DownloadsDirectory + "alvarocontainer1obj1")
	originalHash := md5sum(conf.FilePath)
	if strings.Compare(currentHash, originalHash)!=0 {
		t.Error("Obtained file failed")
	}
	fmt.Println("---------------------------------------------------")


	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlGetAccSuccess.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"400")==0 {
		t.Error("Obj not obtained")

	}else{fmt.Println(" OK")}
	fmt.Println("---------------------------------------------------")

	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlGetAccFailure.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"400")==0 {
		fmt.Println("Failed as expected")
	}else{
		t.Error("Expecting failure")
	}

	fmt.Println("---------------------------------------------------")

	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlGetContSuccess.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"400")!=0 {
		fmt.Println("OK")
	}else{
		t.Error("GET cont failed")
	}
	fmt.Println("---------------------------------------------------")

	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlGetContFailure.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"400")==0 {
		fmt.Println("failed as expected")
	}else{
		t.Error("Expecting failure")
	}

	fmt.Println("---------------------------------------------------")
	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlGetContFailure2.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"400")==0 {
		fmt.Println("failed as expected")
	}else{
		t.Error("Expecting failure")
	}

	fmt.Println("---------------------------------------------------")

	if cmdOut, err = exec.Command(path+"/shellScriptsTests/curlGetContFailure3.sh").Output(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error running command: ", err)
		os.Exit(1)
	}
	resp = string(cmdOut)
	fmt.Println("curl response ", resp)
	if strings.Compare(resp,"400")==0 {
		fmt.Println("failed as expected")
	}else{
		t.Error("Expecting failure")
	}




	time.AfterFunc(600 * time.Second, func(){
		if err:= peer1arun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer2arun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer3arun.Shutdown(nil); err!=nil{
			panic(err)
		}

		if err:= peer1brun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer2brun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer3brun.Shutdown(nil); err!=nil{
			panic(err)
		}

		if err:= peer1crun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer2brun.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= peer3crun.Shutdown(nil); err!=nil{
			panic(err)
		}


		if err:= proxy1run.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= proxy2run.Shutdown(nil); err!=nil{
			panic(err)
		}
		if err:= proxy3run.Shutdown(nil); err!=nil{
			panic(err)
		}

		if err:= tracker.Shutdown(nil); err!=nil{
			panic(err)
		}

	})
	time.Sleep(20*time.Second)
}


