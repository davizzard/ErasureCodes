package httpGo


import (
	"fmt"
	"os"
	"net/http"
	"encoding/json"
)

func CheckSimpleErr(err error, channel chan bool, exit bool) bool {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		if channel != nil {
			channel <- false
		}
		if exit {
			os.Exit(2)
		}
		return true
	} else {
		return false
	}
}

func CheckComplexErr(err error, channel chan bool, fileToRemove string, fileToClose os.File, exit bool) bool {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		if channel != nil {
			channel <- false
		}
		if fileToRemove != "" {
			os.Remove(fileToRemove)
		}
		fileToClose.Close()
		if exit {
			os.Exit(2)
		}
		return true
	} else {
		return false
	}
}

func CheckJsonErr(err error, channel chan bool, w http.ResponseWriter) bool {
	if err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		fmt.Println("Error: %s", err.Error())
		err := json.NewEncoder(w).Encode(err)
		CheckSimpleErr(err, channel, false)
		if channel != nil {
			channel <- false
		}
		return true
	} else {
		return false
	}
}

func CheckLengthErr(lenght int, msg string, channel chan bool, exit bool) bool {
	if lenght == 0 {
		fmt.Println("Error: ", msg)
		if channel != nil {
			channel <- false
		}
		if exit {
			os.Exit(2)
		}
		return true
	} else {
		return false
	}
}
