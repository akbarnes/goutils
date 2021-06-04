package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
	"path/filepath"
	"encoding/json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

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

func (snap Snapshot) Write(snapshotPath string) {
	f, err := os.Create(snapshotPath)

	if err != nil {
		panic(fmt.Sprintf("Error: Could not create snapshot file %s", snapshotPath))
	}
	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(snap)
	f.Close()
}

var TimeStampMode bool
var NumChars int
var Message string

func init() {
	flag.BoolVar(&TimeStampMode, "t", false, "prefix with time stamp")
	flag.IntVar(&NumChars, "n", 40, "number of characters to print")
	flag.StringVar(&Message, "m", "", "commit message")
}

type Snapshot struct {
	ID            string
	Message       string
	Time          string
	Files	      []string
	StoredFiles	  []string
}

func main() {
	flag.Parse()

	// Optional timestamp
	t := time.Now()
	ts := t.Format("2006-01-02T15-04-05")
	snap := Snapshot{Time: ts}
	snap.Files = []string{}
	snap.StoredFiles = []string{}
	snap.Message = Message

	for i := 0; i < flag.NArg(); i++ {
		fileName := flag.Arg(i)
		ext := filepath.Ext(fileName)
		hash, _ := HashFile(fileName, NumChars)
		verFolder := filepath.Join(".gover", "data", hash[0:2]) 
		verFile := filepath.Join(verFolder, hash + ext)

		snap.Files = append(snap.Files, fileName)
		snap.StoredFiles = append(snap.StoredFiles, verFile)

		os.MkdirAll(verFolder, 0777)
		CopyFile(fileName, verFile)
		fmt.Printf("%s -> %s\n", fileName, verFile)
	}

	snapFolder := filepath.Join(".gover", "snapshots")
	os.MkdirAll(snapFolder, 0777)
	snapFile := filepath.Join(snapFolder, ts + ".json")
	snap.Write(snapFile)
}

