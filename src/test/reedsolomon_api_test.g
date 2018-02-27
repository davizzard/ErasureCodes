package test

import (
	"davizzard/ErasureCodes/src/API"
	"testing"
	"fmt"
)

const fileChunk = 1*(1<<10) // 1 KB

func TestEncode(t *testing.T) {
	fmt.Println("Begin TEST ENCODE --------------------------------")
	API.EncodeFileAPI("/root/Desktop/empty/go/src/davizzard/ErasureCodes/src/test/test.txt", fileChunk, 3, nil)
	fmt.Println("End TEST ENCODE --------------------------------")
}

func TestDecode(t *testing.T) {
	fmt.Println("Begin TEST DECODE --------------------------------")
	//API.DecodeFileAPI("/root/Desktop/empty/go/src/davizzard/ErasureCodes/src/test/test.txt", 40, 40, 3, ".")
	fmt.Println("End TEST DECODE --------------------------------")
}
