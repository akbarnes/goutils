package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	// "os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func HashFile(FileName string, NumChars int, TimeStampMode bool) (string, error) {
	var data []byte
	var err error

	// Optional timestamp
	t := time.Now()

	if TimeStampMode {
		fmt.Print(t.Format("2006-01-02T15-04-05-"))
	}


	data, err = ioutil.ReadFile(FileName)

	if err != nil {
		return "", err
	}


	sum := fmt.Sprintf("%x", sha256.Sum256(data))

	if len(sum) < NumChars || NumChars < 0 {
		NumChars = len(sum)
	}

	return sum[0:NumChars], nil
}

var TimeStampMode bool
var NumChars int

func init() {
	flag.BoolVar(&TimeStampMode, "t", false, "prefix with time stamp")
	flag.IntVar(&NumChars, "n", 32, "number of characters to print")
}

func main() {
	flag.Parse()

		for i := 0; i < flag.NArg(); i++ {
			hash, _ := HashFile(flag.Arg(i), NumChars, TimeStampMode)
			fmt.Printf("%s -> %s\n", flag.Arg(i), hash)
		}
}
