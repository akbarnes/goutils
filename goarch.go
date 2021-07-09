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
	Lengths  map[string]int64
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

	defer f.Close()
	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(snap)
}

func StoreFolder(archivePrefix string, workingDirectory string) (int64, error) {
	snap := Snapshot{}
	snap.Files = []string{}
	snap.Offsets = make(map[string]int64)
	snap.Lengths = make(map[string]int64)
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
		snap.Lengths[fileName] = fileBytes
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
// 	Check(err)

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

func ListArchiveContents(archivePrefix string) {
	snap, err := ReadArchive(archivePrefix)

	if err != nil {
		fmt.Println("Error reading archive contents")
	}

	for i, file := range snap.Files {
		mtime := snap.ModTimes[file]
		bytes := snap.Lengths[file]

		if VerboseMode {
			fmt.Printf("%03d: %19s, %4d MB, %s\n", i, mtime, bytes/1000000, file)
		} else {
			fmt.Println(file)
		}
	}
}

// Read a snapshot given a file path
func ReadArchive(archivePrefix string) (Snapshot, error) {
	archivePath := archivePrefix + ".json"

	var mySnapshot Snapshot
	f, err := os.Open(archivePath)

	if err != nil {
		// panic(fmt.Sprintf("Error: Could not read snapshot file %s", snapshotPath))
		return Snapshot{}, err
	}

	defer f.Close()
	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&mySnapshot); err != nil {
		fmt.Printf("Error:could not decode archive file %s\n", archivePath)
		Check(err)
	}

	return mySnapshot, nil
}

var VerboseMode bool

func AddOptionFlags(fs *flag.FlagSet) {
	fs.BoolVar(&VerboseMode, "verbose", false, "verbose mode")
	fs.BoolVar(&VerboseMode, "v", false, "verbose mode")
}

func main() {
	storeCmd := flag.NewFlagSet("store", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	// extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Expected subcommand")
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "store" || cmd == "st" || cmd == "s" {
		AddOptionFlags(storeCmd)
		storeCmd.Parse(os.Args[2:])
		archivePath := storeCmd.Arg(0)
		inputFolder := storeCmd.Arg(1)
		StoreFolder(archivePath, inputFolder)
	} else if cmd == "list" || cmd == "ls" || cmd == "l" {
		AddOptionFlags(listCmd)
		listCmd.Parse(os.Args[2:])
		archivePrefix := listCmd.Arg(0)
		ListArchiveContents(archivePrefix)
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
