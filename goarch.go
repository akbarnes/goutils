package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Snapshot struct {
	Files    []string
	ModTimes map[string]string
	Offsets  map[string]int64
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// Copy the source file to a destination file. Any existing file
// will be overwritten and will not copy file attributes.
func CopyFile(out *os.File, src string) (int64, error) {
	in, err := os.Open(src)

	if err != nil {
		return 0, err
	}

	defer in.Close()

	nbytes, err := io.Copy(out, in)

	if err != nil {
		return 0, err
	}

	return nbytes, nil
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

func StoreFolder(archivePrefix string, workingDirectory string) (int64, error) {
	snap := Snapshot{}
	snap.Files = []string{}
	snap.Offsets = make(map[string]int64)
	snap.ModTimes = make(map[string]string)

	archivePath := archivePrefix + ".dat"
	archiveFile, err := os.Create(archivePath)
	var nbytes int64 = 0

	if err != nil {
		return nbytes, err
	}

	defer archiveFile.Close()

	var VersionFile = func(fileName string, info os.FileInfo, err error) error {
		fileName = strings.TrimSuffix(fileName, "\n")

		if info.IsDir() {
			return nil
		}

		props, err := os.Stat(fileName)

		if err != nil {
			if VerboseMode {
				fmt.Printf("Skipping unreadable file %s\n", fileName)
			}

			return nil
		}

		fmt.Println(fileName)
		modTime := props.ModTime().Format("2006-01-02T15-04-05")
		snap.Files = append(snap.Files, fileName)
		snap.ModTimes[fileName] = modTime
		snap.Offsets[fileName] = nbytes
		fileBytes, _ := CopyFile(archiveFile, fileName)
		nbytes += fileBytes
		return nil
	}

	// fmt.Printf("No changes detected in %s for commit %s\n", workDir, snapshot.ID)
	filepath.Walk(workingDirectory, VersionFile)
	snapFile := archivePrefix + ".json"
	snap.Write(snapFile)
	return nbytes, nil
}

// func CheckoutSnaphot(snapshotNum int, outputFolder string) {
// 	if len(outputFolder) == 0 {
// 		outputFolder = fmt.Sprintf("snapshot%04d", snapshotNum)
// 	}

// 	fmt.Printf("Checking out %s\n", snapshotNum)

// 	snapshotGlob := filepath.Join(".gover", "snapshots", "*.json")
// 	snapshotPaths, err := filepath.Glob(snapshotGlob)
// 	check(err)

// 	snapshotPath := snapshotPaths[snapshotNum-1]
// 	fmt.Printf("Reading %s\n", snapshotPath)
// 	snap := ReadSnapshotFile(snapshotPath)

// 	os.Mkdir(outputFolder, 0777)

// 	for _, file := range snap.Files {
// 		fileDir := filepath.Dir(file)
// 		outDir := outputFolder

// 		if fileDir != "." {
// 			outDir = filepath.Join(outputFolder, fileDir)
// 			fmt.Printf("Creating folder %s\n", outDir)
// 			os.MkdirAll(outDir, 0777)
// 		}

// 		outFile := filepath.Join(outputFolder, file)
// 		storedFile := snap.StoredFiles[file]
// 		fmt.Printf("Restoring %s to %s\n", storedFile, outFile)
// 		CopyFile(storedFile, outFile)
// 	}
// }

// func LogSingleSnapshot(snapshotNum int) {
// 	snapshotGlob := filepath.Join(".gover", "snapshots", "*.json")
// 	snapshotPaths, err := filepath.Glob(snapshotGlob)
// 	check(err)

// 	snapshotPath := snapshotPaths[snapshotNum-1]

// 	snap := ReadSnapshotFile(snapshotPath)

// 	if JsonMode {
// 		type SnapshotFile struct {
// 			File       string
// 			StoredFile string
// 		}

// 		snapFiles := []SnapshotFile{}

// 		for _, file := range snap.Files {
// 			snapFile := SnapshotFile{File: file, StoredFile: snap.StoredFiles[file]}
// 			snapFiles = append(snapFiles, snapFile)
// 		}

// 		PrintJson(snapFiles)
// 	} else {
// 		for _, file := range snap.Files {
// 			fmt.Println(file)
// 		}
// 	}
// }

// func ReadSnapshot(snapId string) Snapshot {
// 	snapshotPath := filepath.Join(".gover", "snapshots", snapId+".json")

// 	if VerboseMode {
// 		fmt.Printf("Reading %s\n", snapshotPath)
// 	}

// 	return ReadSnapshotFile(snapId)
// }

// // Read a snapshot given a file path
// func ReadSnapshotFile(snapshotPath string) Snapshot {
// 	var mySnapshot Snapshot
// 	f, err := os.Open(snapshotPath)

// 	if err != nil {
// 		// panic(fmt.Sprintf("Error: Could not read snapshot file %s", snapshotPath))
// 		return Snapshot{Files: []string{}, StoredFiles: make(map[string]string), ModTimes: make(map[string]string)}
// 	}

// 	myDecoder := json.NewDecoder(f)

// 	if err := myDecoder.Decode(&mySnapshot); err != nil {
// 		fmt.Printf("Error:could not decode head file %s\n", snapshotPath)
// 		check(err)
// 	}

// 	f.Close()
// 	return mySnapshot
// }

var VerboseMode bool

func init() {
	flag.BoolVar(&VerboseMode, "verbose", false, "verbose")
	flag.BoolVar(&VerboseMode, "v", false, "verbose")
}

// type Commit struct {
// 	ID        string
// 	Branch    string
// 	Message   string
// 	Time      string
// 	ParentIDs []string
// 	Files     []fileInfo
// 	ChunkIDs  []string
// }

func main() {
	storeCmd := flag.NewFlagSet("store", flag.ExitOnError)
	// listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	// extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Expected subcommand")
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "store" || cmd == "st" || cmd == "s" {
		storeCmd.Parse(os.Args[2:])
		archivePath := storeCmd.Arg(0)
		inputFolder := storeCmd.Arg(1)
		StoreFolder(archivePath, inputFolder)
		// } else if cmd == "list" || cmd == "ls" || cmd == "l" {
		// 	listCmd.Parse(os.Args[2:])
		// 	archiveFile = commitCmd.Arg(0)
		// 	ListArchiveContents(archiveFile)
		// } else if cmd == "extract" || cmd == "ex" || cmd == "e" || cmd == "x" {
		// 	extractCmd.Parse(os.Args[2:])
		// 	archiveFile = commitCmd.Arg(0)
		// 	outputFolder = commitCmd.Arg(1)
		// 	ExtractArchive(archiveFile, outputFolder)
	} else {
		fmt.Println("Unknown subcommand")
		os.Exit(1)
	}
}
