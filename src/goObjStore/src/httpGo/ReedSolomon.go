package httpGo

import (
	"fmt"
	"os"
	"github.com/klauspost/reedsolomon"
	"path/filepath"
	"io"
	"strconv"
	"math"
	"time"
	"github.com/davizzard/ErasureCodes/src/goObjStore/src/conf"
	"sync"
)


func EncodeFileAPI(fname string, fileChunk int, parityShards int, putOK chan bool) (int, string, int) {
	defer elapsed("EncodeFileAPI")()
	var dataShards int
	var counter int = 0

	fmt.Println("Opening", fname)
	f, err := os.Open(fname)
	CheckSimpleErr(err, putOK, true)

	fInfo, err := f.Stat()
	CheckSimpleErr(err, putOK, true)


	text := strconv.FormatInt(fInfo.Size(), 10)        // size
	size, err := strconv.Atoi(text)
	CheckComplexErr(err, putOK, fname, *f, true)

	fmt.Printf("Filechunk size BEFORE: %d bytes.\n", fileChunk)

	dataShards = int(math.Ceil(float64(size) / float64(fileChunk)))
	if (dataShards > 257) {
		dataShards = 257
	}

	fmt.Println("Finding the correct amount of Data shards...")
	for math.Ceil(math.Mod(float64(size), float64(dataShards))) != 0 {
		fmt.Printf("Iteration nº: %d\n", counter)
		dataShards--
		counter++
	}

	fileChunk = int(math.Ceil(float64(size) / float64(dataShards+parityShards)))

	fmt.Printf("Filechunk size AFTER: %d bytes.\n", fileChunk)

	// Create encoding matrix.
	enc, err := reedsolomon.NewStreamC(dataShards, parityShards, true, true)
	CheckSimpleErr(err, putOK, true)
	shards := dataShards + parityShards
	out := make([]*os.File, shards)

	// Create the resulting files.
	_, file := filepath.Split(fname)
	dir := conf.LocalDirectory + "/" + file
	os.Mkdir(dir, 0777)
	dir = dir + "/"
	for i := range out {
		outfn := fmt.Sprintf("%d", i)
		fmt.Println("Creating", outfn)
		out[i], err = os.Create(filepath.Join(dir, outfn))
		CheckSimpleErr(err, putOK, true)

	}

	// Split into files.
	data := make([]io.Writer, dataShards)
	for i := range data {
		data[i] = out[i]
	}
	// Do the split
	err = enc.Split(f, data, fInfo.Size())
	CheckSimpleErr(err, putOK, true)

	// Close and re-open the files.
	input := make([]io.Reader, dataShards)

	for i := range data {
		out[i].Close()
		f, err := os.Open(out[i].Name())
		CheckSimpleErr(err, putOK, true)
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
	CheckSimpleErr(err, putOK, true)
	fmt.Printf("File split into %d data + %d parity shards.\n", dataShards, parityShards)

	return shards, file, size

}



func DecodeFileAPI(fname string, key string, dataShards int, parityShards int, separator string, nodeList [][]string) (bool) {
	defer elapsed("DecodeFileAPI")()
	var missingShards []int
	var exitStatus = false

	fmt.Printf("Decoding file... Data Shards: %d. Parity shards: %d\n", dataShards, parityShards)
	// Create matrix
	enc, err := reedsolomon.NewStreamC(dataShards, parityShards, true, true)
	CheckSimpleErr(err, nil, true)

	// Open the inputs
	shards, size, err := openInput(dataShards, parityShards, fname, separator)
	CheckSimpleErr(err, nil, true)

	// Verify the shards
	ok, err := enc.Verify(shards)
	//ok := true
	if ok {
		fmt.Println("No reconstruction needed")
	} else {
		fmt.Println("Verification failed. Reconstructing data.")

		shards, size, err = openInput(dataShards, parityShards, fname, separator)
		CheckSimpleErr(err, nil, true)
		// Create out destination writers
		out := make([]io.Writer, len(shards))
		for i := range out {
			if shards[i] == nil {
				outfn := fmt.Sprintf("%s%s%d", fname, separator, i)
				fmt.Println("Creating", outfn)
				out[i], err = os.Create(outfn)
				CheckSimpleErr(err, nil, true)
				missingShards = append(missingShards, i)
			}
		}
		err = enc.Reconstruct(shards, out)
		if err != nil {
			CheckSimpleErr(err, nil, false)
			exitStatus = true
			return exitStatus
		}
		fmt.Printf("Missing shards: %d\n", len(missingShards))

		//Sending reconstructed shards to Nodes
		var wg sync.WaitGroup
		wg.Add(len(missingShards))

		fmt.Println("Resending missing shards to Storage Nodes...")
		for _, j := range missingShards {
			nodeNum := j % len(nodeList)
			fileShardPath := fname + separator + strconv.Itoa(j)

			SendFileToNodes(fileShardPath, nodeList, key, nodeNum, 0, j, &wg, "")

			fmt.Printf("Data shard %d reconstructed. Sending data to Storage Node with ID %d.\n", j, nodeNum)

		}

		wg.Wait()

		// Close output.
		for i := range out {
			if out[i] != nil {
				err := out[i].(*os.File).Close()
				CheckSimpleErr(err, nil, true)
			}
		}
		shards, size, err = openInput(dataShards, parityShards, fname, separator)
		ok, err = enc.Verify(shards)
		if !ok {
			fmt.Println("Verification failed after reconstruction, data likely corrupted:", err)
			exitStatus = true
			return exitStatus
		}
		CheckSimpleErr(err, nil, true)

	}

	// Join the shards and write them
	outfn := conf.DownloadsDirectory + key

	fmt.Println("Writing data to", outfn)
	f, err := os.Create(outfn)
	CheckSimpleErr(err, nil, true)

	shards, size, err = openInput(dataShards, parityShards, fname, separator)
	CheckSimpleErr(err, nil, true)

	err = enc.Join(f, shards, int64(dataShards)*size)
	CheckSimpleErr(err, nil, true)

	return exitStatus
}


func openInput(dataShards, parShards int, fname string, separator string) (r []io.Reader, size int64, err error) {
	// Create shards and load the data.
	shards := make([]io.Reader, dataShards+parShards)
	for i := range shards {
		infn := fmt.Sprintf("%s%s%d", fname, separator, i)
		//fmt.Println("Opening", infn)
		f, err := os.Open(infn)
		if err != nil {
			fmt.Println("Error reading file", err)
			shards[i] = nil
			continue
		} else {
			shards[i] = f
		}
		stat, err := f.Stat()
		CheckSimpleErr(err, nil, true)
		if stat.Size() > 0 {
			size = stat.Size()
		} else {
			shards[i] = nil
		}
	}
	return shards, size, nil
}


func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}