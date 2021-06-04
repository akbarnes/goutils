package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	// "time"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// // Optional timestamp
// t := time.Now()

// if TimeStampMode {
// 	fmt.Print(t.Format("2006-01-02T15-04-05-"))
// }

func HashFile(FileName string, NumChars int) (string, error) {
	var data []byte
	var err error

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

// Copy the source file to a destination file. Any existing file
// will be overwritten and will not copy file attributes.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

var TimeStampMode bool
var NumChars int

func init() {
	flag.BoolVar(&TimeStampMode, "t", false, "prefix with time stamp")
	flag.IntVar(&NumChars, "n", 64, "number of characters to print")
}

func main() {
	os.Mkdir(".gover", 0777)

	flag.Parse()

		for i := 0; i < flag.NArg(); i++ {
			fileName := flag.Arg(i)
			ext := filepath.Ext(fileName)
			hash, _ := HashFile(fileName, NumChars)

			verFolder := filepath.Join(".gover", "data", hash[0:2]) 
			os.MkdirAll(verFolder, 0777)

			verFile := filepath.Join(verFolder, hash + ext)
			fmt.Printf("%s -> %s\n", fileName, verFile)
		}
}
