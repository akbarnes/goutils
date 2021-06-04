package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
	"math/rand"
	"path/filepath"
	"encoding/json"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))
  
const HexChars = "0123456789abcdef"

// Return a random string of specified length with hexadecimal characters
func RandHexString(length int) string {
	return RandString(length, HexChars)
}

// Return a random string of specified length with an arbitrary character set
func RandString(length int, charset string) string {
	b := make([]byte, length)

	for i := range b {
	b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
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

// Read a snapshot given a file path
func ReadSnapshotFile(snapshotPath string) Snapshot {
	var mySnapshot Snapshot
	f, err := os.Open(snapshotPath)

	if err != nil {
		// panic(fmt.Sprintf("Error: Could not read snapshot file %s", snapshotPath))
		return Snapshot{}
	}

	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&mySnapshot); err != nil {
		panic(fmt.Sprintf("Error:could not decode snapshot file %s", snapshotPath))
	}

	f.Close()
	return mySnapshot
}

var LogCommand bool
var LogFiles bool
var NumChars int
var Message string

func init() {
	flag.BoolVar(&LogCommand, "log", false, "list snapshots")
	flag.BoolVar(&LogFiles, "f", false, "print files when listing snapshots")
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

	if LogCommand {
		snapshotGlob := filepath.Join(".gover","snapshots","*.json")
		snapshotPaths, err := filepath.Glob(snapshotGlob)
		check(err)

		for _, snapshotPath := range snapshotPaths {
			snap := ReadSnapshotFile(snapshotPath)

			// ID: 943e8daa (943e8daa4bc0ab899c36b5030d4a27a6b833b2ba)
			// Time: 2021/05/08 08:57:46
			// Message: specify workdir path explicitly
			fmt.Printf("Time: %s\nID: %s\nMessage: %s\n", snap.Time, snap.ID, snap.Message)

			if LogFiles {
				fmt.Println("Files:")

				for _, file := range snap.Files {
					fmt.Printf("  %s\n",file)
				}
			}
			
			fmt.Println("")
		}
	} else {
		// Optional timestamp
		t := time.Now()
		ts := t.Format("2006-01-02T15-04-05")
		snap := Snapshot{Time: ts}
		snap.Files = []string{}
		snap.StoredFiles = []string{}
		snap.ID = RandHexString(40)
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
}

