package API

import (
	"fmt"
	"os"
	"github.com/klauspost/reedsolomon"
	"path/filepath"
	"io"
	"strconv"
	"davizzard/ErasureCodes/src/goObjStore/src/httpVar"
	"math"
	"davizzard/ErasureCodes/src/goObjStore/src/conf"
	"time"
)


func EncodeFileAPI(fname string, fileChunk int, parityShards int, putOK chan bool) (int, string, int) {
	defer elapsed("EncodeFileAPI")()
	var dataShards int
	var counter int = 0

	fmt.Println("Opening", fname)
	f, err := os.Open(fname)
	CheckErr(err)

	fInfo, err := f.Stat()
	CheckErr(err)


	text := strconv.FormatInt(fInfo.Size(), 10)        // size
	size, _ := strconv.Atoi(text)
	if err != nil {
		fmt.Println(err.Error())
		os.Remove(fname)
		f.Close()
		putOK <- false
		CheckErr(err)
	}
	/*
	for math.Ceil(math.Mod(float64(size), float64(fileChunk))) != 0 {
		fmt.Print("Iteration nº: ")
		fmt.Println(counter)
		fileChunk++
		counter++
	}
	*/
	fmt.Print("Filechunk size BEFORE is: ")
	fmt.Println(fileChunk)

	httpVar.TotalNumMutex.Lock()
	dataShards = int(math.Ceil(float64(size) / float64(fileChunk)))
	httpVar.TotalNumMutex.Unlock()

	for math.Ceil(math.Mod(float64(size), float64(dataShards))) != 0 {
		fmt.Print("Iteration nº: ")
		fmt.Println(counter)
		dataShards++
		counter++
	}

	fileChunk = int(math.Ceil(float64(size) / float64(dataShards+parityShards)))

	fmt.Print("Filechunk size AFTER is: ")
	fmt.Println(fileChunk)

	// Checking arguments.
	if dataShards > 257 {
		fmt.Fprintf(os.Stderr, "Error: Too many data shards\n")
		os.Exit(1)
	}

	// Create encoding matrix.
	enc, err := reedsolomon.NewStream(dataShards, parityShards)
	CheckErr(err)


	shards := dataShards + parityShards
	out := make([]*os.File, shards)

	// Create the resulting files.
	_, file := filepath.Split(fname)
	dir := conf.LocalDirectory

	for i := range out {
		outfn := fmt.Sprintf("%s.%d", file, i)
		//outfn := fmt.Sprintf("NEW%d", i)
		fmt.Println("Creating", outfn)
		out[i], err = os.Create(filepath.Join(dir, outfn))
		CheckErr(err)
	}

	// Split into files.
	data := make([]io.Writer, dataShards)
	for i := range data {
		data[i] = out[i]
	}
	// Do the split
	err = enc.Split(f, data, fInfo.Size())
	CheckErr(err)

	// Close and re-open the files.
	input := make([]io.Reader, dataShards)

	for i := range data {
		out[i].Close()
		f, err := os.Open(out[i].Name())
		CheckErr(err)
		input[i] = f
		defer f.Close()
	}

	// Create parity output writers
	parity := make([]io.Writer, parityShards)
	for i := range parity {
		parity[i] = out[dataShards+i]
		defer out[dataShards+i].Close()
	}

	// Encode parity
	err = enc.Encode(input, parity)
	CheckErr(err)
	fmt.Printf("File split into %d data + %d parity shards.\n", dataShards, parityShards)

	return shards, fInfo.Name(), size

}



func DecodeFileAPI(fname string, key string, dataShards int, parityShards int, separator string, putOK chan bool) {
	defer elapsed("DecodeFileAPI")()
	fmt.Print("Decoding file... Data Shards: ")
	fmt.Print(dataShards)
	fmt.Print(". Parity shards: ")
	fmt.Println(parityShards)
	// Create matrix
	enc, err := reedsolomon.NewStream(dataShards, parityShards)
	CheckErr(err)

	// Open the inputs
	shards, size, err := openInput(dataShards, parityShards, fname, separator)
	CheckErr(err)

	// Verify the shards
	ok, err := enc.Verify(shards)
	//ok := true
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Verification failed. Reconstructing data.")

		shards, size, err = openInput(dataShards, parityShards, fname, separator)
		CheckErr(err)
		// Create out destination writers
		out := make([]io.Writer, len(shards))
		for i := range out {
			if shards[i] == nil {
				outfn := fmt.Sprintf("%s%s%d", fname, separator, i)
				fmt.Println("Creating", outfn)
				out[i], err = os.Create(outfn)
				CheckErr(err)
			}
		}
		err = enc.Reconstruct(shards, out)
		if err != nil {
			fmt.Println("Reconstruct failed -", err)
			os.Exit(1)
		}
		// Close output.
		for i := range out {
			if out[i] != nil {
				err := out[i].(*os.File).Close()
				CheckErr(err)
			}
		}
		shards, size, err = openInput(dataShards, parityShards, fname, separator)
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted:", err)
			os.Exit(1)
		}
		CheckErr(err)

	}

	// Join the shards and write them
	outfn := conf.DownloadsDirectory + key

	fmt.Println("Writing data to", outfn)
	f, err := os.Create(outfn)
	CheckErr(err)

	shards, size, err = openInput(dataShards, parityShards, fname, separator)
	CheckErr(err)

	// We don't know the exact filesize.
	err = enc.Join(f, shards, int64(dataShards)*size)
	CheckErr(err)

}


func openInput(dataShards, parShards int, fname string, separator string) (r []io.Reader, size int64, err error) {
	// Create shards and load the data.
	shards := make([]io.Reader, dataShards+parShards)
	for i := range shards {
		infn := fmt.Sprintf("%s%s%d", fname, separator, i)
		fmt.Println("Opening", infn)
		f, err := os.Open(infn)
		if err != nil {
			fmt.Println("Error reading file", err)
			shards[i] = nil
			continue
		} else {
			shards[i] = f
		}
		stat, err := f.Stat()
		CheckErr(err)
		if stat.Size() > 0 {
			size = stat.Size()
		} else {
			shards[i] = nil
		}
	}
	return shards, size, nil
}



func CheckErr(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(2)
	}
}


func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}