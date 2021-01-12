package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var TimeStampMode bool
var NumChars int

func init() {
	flag.BoolVar(&TimeStampMode, "t", false, "prefix with time stamp")
	flag.IntVar(&NumChars, "n", 64, "number of characters to print")

}

func main() {
	flag.Parse()

	var data []byte
	var err error

	// Optional timestamp
	t := time.Now()

	if TimeStampMode {
		fmt.Print(t.Format("2006-01-02T15-04-05-"))
	}

	if flag.NArg() >= 1 {
		FileName := flag.Arg(1)
		data, err = ioutil.ReadFile(FileName)

		if err != nil {
			panic("Could not read file " + FileName)
		}
	} else {
		data, err = ioutil.ReadAll(os.Stdin)
		check(err)
	}

	sum := fmt.Sprintf("%x", sha256.Sum256(data))

	if len(sum) < NumChars || NumChars < 0 {
		NumChars = len(sum)
	}

	fmt.Printf("%s\n", sum[0:NumChars])
}
