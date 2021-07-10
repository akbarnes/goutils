package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

// Copy the source file to a destination file. Any existing file
// will be overwritten and will not copy file attributes.
func SplitFile(outPrefix string, src string, maxBytes int64) (int, error) {
	props, err := os.Stat(src)

	bytesRemaining := props.Size()
	fileCount := 0

	outputFolder := filepath.Dir(outPrefix)

	if outputFolder != "." {
		if err := os.MkdirAll(outPrefix, 0777); err != nil {
			if VerboseMode {
				fmt.Printf("Error creating output folder %s\n", outputFolder)
			}

			return 0, err
		}
	}

	in, err := os.Open(src)

	if err != nil {
		return fileCount, err
	}

	defer in.Close()

	for bytesRemaining > 0 {
		fileCount++

		outPath := fmt.Sprintf("%s.%d.part", outPrefix, fileCount)
		out, err := os.Create(outPath)

		if err != nil {
			return 0, err
		}

		defer out.Close()

		copyBytes := min64(maxBytes, bytesRemaining)
		bytesCopied, err := io.CopyN(out, in, copyBytes)

		if err == nil {
			bytesRemaining -= bytesCopied
		} else {
			return fileCount, err
		}

		if VerboseMode {
			fmt.Printf("%03d: %4d MB, %s\n", fileCount, bytesCopied/(1000*1000), outPath)
		} else {
			fmt.Println(outPath)
		}

	}

	return fileCount, nil
}

// func ExtractFile(file string, snap Snapshot, in *os.File, dst string) error {
// 	outputFolder := filepath.Dir(dst)

// 	if err := os.MkdirAll(outputFolder, 0777); err != nil {
// 		if VerboseMode {
// 			fmt.Printf("Error creating output folder %s\n", outputFolder)
// 		}

// 		return err
// 	}

// 	out, err := os.Create(dst)

// 	if err != nil {
// 		if VerboseMode {
// 			fmt.Printf("Error creating destination file %s:\n", dst)
// 			fmt.Println(err)
// 			fmt.Printf("\n")
// 		}

// 		return err
// 	}

// 	defer out.Close()

// 	offset := snap.Offsets[file]
// 	nbytes := snap.Lengths[file]

// 	if _, err := in.Seek(offset, 0); err != nil {
// 		if VerboseMode {
// 			fmt.Println("Error seeking on archive file")
// 		}

// 		return err
// 	}

// 	if _, err := io.CopyN(out, in, nbytes); err != nil {
// 		if VerboseMode {
// 			fmt.Printf("Error copying to destination file %s\n", dst)
// 			fmt.Println(err)
// 			fmt.Printf("\n")
// 		}

// 		return err
// 	}

// 	return nil
// }

// func (snap Snapshot) Write(snapshotPath string) {
// 	f, err := os.Create(snapshotPath)

// 	if err != nil {
// 		panic(fmt.Sprintf("Error: Could not create snapshot file %s", snapshotPath))
// 	}

// 	defer f.Close()
// 	myEncoder := json.NewEncoder(f)
// 	myEncoder.SetIndent("", "  ")
// 	myEncoder.Encode(snap)
// }

// func StoreFolder(archivePrefix string, workingDirectory string) {
// 	snap := Snapshot{}
// 	snap.Files = []string{}
// 	snap.Offsets = make(map[string]int64)
// 	snap.Lengths = make(map[string]int64)
// 	snap.ModTimes = make(map[string]string)

// 	archivePath := archivePrefix + ".dat"
// 	archiveFile, err := os.Create(archivePath)
// 	var nbytes int64 = 0

// 	if err != nil {
// 		fmt.Printf("Error creating archive file %s\n", archivePath)
// 		return
// 	}

// 	defer archiveFile.Close()

// 	var VersionFile = func(fileName string, info os.FileInfo, err error) error {
// 		fileName = strings.TrimSuffix(fileName, "\n")

// 		if info.IsDir() {
// 			return nil
// 		}

// 		props, err := os.Stat(fileName)

// 		if err != nil {
// 			if VerboseMode {
// 				fmt.Printf("Skipping unreadable file %s\n", fileName)
// 			}

// 			return nil
// 		}

// 		fmt.Println(fileName)
// 		modTime := props.ModTime().Format("2006-01-02T15-04-05")
// 		snap.Files = append(snap.Files, fileName)
// 		snap.ModTimes[fileName] = modTime
// 		snap.Offsets[fileName] = nbytes
// 		fileBytes, _ := StoreFile(archiveFile, fileName)
// 		snap.Lengths[fileName] = fileBytes
// 		nbytes += fileBytes
// 		return nil
// 	}

// 	// fmt.Printf("No changes detected in %s for commit %s\n", workDir, snapshot.ID)
// 	filepath.Walk(workingDirectory, VersionFile)
// 	snapFile := archivePrefix + ".json"
// 	snap.Write(snapFile)
// }

// func ExtractArchive(archivePrefix string, outputFolder string) {
// 	snap, err := ReadArchive(archivePrefix)

// 	if err != nil {
// 		fmt.Println("Error reading archive contents")
// 	}

// 	archivePath := archivePrefix + ".dat"
// 	archiveFile, err := os.Open(archivePath)

// 	if err != nil {
// 		fmt.Printf("Cannot open archive file %s\n", archivePath)
// 		return
// 	}

// 	defer archiveFile.Close()

// 	for _, file := range snap.Files {
// 		outPath := filepath.Join(outputFolder, file)

// 		if err := ExtractFile(file, snap, archiveFile, outPath); err == nil {
// 			fmt.Println(file)
// 		}
// 	}
// }

// // Read a snapshot given a file path
// func ReadArchive(archivePrefix string) (Snapshot, error) {
// 	archivePath := archivePrefix + ".json"

// 	var mySnapshot Snapshot
// 	f, err := os.Open(archivePath)

// 	if err != nil {
// 		// panic(fmt.Sprintf("Error: Could not read snapshot file %s", snapshotPath))
// 		return Snapshot{}, err
// 	}

// 	defer f.Close()
// 	myDecoder := json.NewDecoder(f)

// 	if err := myDecoder.Decode(&mySnapshot); err != nil {
// 		fmt.Printf("Error:could not decode archive file %s\n", archivePath)
// 		Check(err)
// 	}

// 	return mySnapshot, nil
// }

var VerboseMode bool
var Output string

func AddOptionFlags(fs *flag.FlagSet) {
	fs.BoolVar(&VerboseMode, "verbose", false, "verbose mode")
	fs.BoolVar(&VerboseMode, "v", false, "verbose mode")
	fs.StringVar(&Output, "out", "", "output path")
	fs.StringVar(&Output, "o", "", "output path")
}

func main() {
	splitCmd := flag.NewFlagSet("split", flag.ExitOnError)
	joinCmd := flag.NewFlagSet("join", flag.ExitOnError)

	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("Expected subcommand")
		os.Exit(1)
	}

	cmd := os.Args[1]

	if cmd == "split" || cmd == "sp" || cmd == "s" {
		AddOptionFlags(splitCmd)
		splitCmd.Parse(os.Args[2:])
		inputPath := splitCmd.Arg(0)
		outputPrefix := Output

		if len(outputPrefix) == 0 {
			outputPrefix = inputPath
		}

		SplitFile(outputPrefix, inputPath, 10*1000*1000)
	} else if cmd == "join" || cmd == "j" {
		AddOptionFlags(joinCmd)
		joinCmd.Parse(os.Args[2:])
		inputPrefix := joinCmd.Arg(0)
		outputPath := Output

		if len(outputPath) == 0 {
			outputPath = inputPrefix
		}

		// JoinFile(outputPath, inputPrefix)
	} else {
		fmt.Println("Unknown subcommand")
		os.Exit(1)
	}
}
