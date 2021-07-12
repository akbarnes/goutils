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
	Files       []string
	PackNumbers map[string][]int
	Offsets     map[string][]int64
	Lengths     map[string][]int64
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

// Copy the source file to a destination file. Any existing file
// will be overwritten and will not copy file attributes.
func StoreFile(out *os.File, src string) (int64, error) {
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

func (snap Snapshot) Write(archiveFolder string) error {
	snapshotPath := filepath.Join(archiveFolder, "snapshot.json")
	return snap.WriteFile(snapshotPath)
}

func (snap Snapshot) WriteFile(snapshotPath string) error {
	f, err := os.Create(snapshotPath)

	if err != nil {
		if VerboseMode {
			fmt.Printf("Error: Could not create snapshot file %s", snapshotPath)
		}

		return err
	}

	defer f.Close()
	myEncoder := json.NewEncoder(f)
	myEncoder.SetIndent("", "  ")
	myEncoder.Encode(snap)
	return nil
}

func min64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func StoreFolder(archiveFolder string, workingDirectory string, maxPackBytes int64) error {
	snap := Snapshot{}
	snap.Files = []string{}
	snap.PackNumbers = make(map[string][]int)
	snap.Offsets = make(map[string][]int64)
	snap.Lengths = make(map[string][]int64)

	if err := os.MkdirAll(archiveFolder, 0777); err != nil {
		if VerboseMode {
			fmt.Printf("Error creating archive folder %s\n", archiveFolder)
		}

		return err
	}

	packCount := 1
	packPath := filepath.Join(archiveFolder, fmt.Sprintf("pack%d.dat", packCount))
	packFile, err := os.Create(packPath)
	var packOffset int64 = 0
	var packBytesRemaining int64 = maxPackBytes

	if err != nil {
		if VerboseMode {
			fmt.Printf("Error creating pack file %s\n", packPath)
		}

		return err
	}

	var VersionFile = func(fileName string, info os.FileInfo, err error) error {
		fileName = strings.TrimSuffix(fileName, "\n")

		if info.IsDir() {
			return nil
		}

		props, err := os.Stat(fileName)

		if err != nil {
			if VerboseMode {
				fmt.Printf("Can't stat file %s, skipping\n", fileName)
			}

			return err
		}

		in, err := os.Open(fileName)

		if err != nil {
			if VerboseMode {
				fmt.Printf("Can't open file %s for reading, skipping\n", fileName)
			}

			return err
		}

		defer in.Close()

		if VerboseMode {
			fmt.Printf("Storing %s\n", fileName)
		} else {
			fmt.Println(fileName)
		}

		fileBytesRemaining := props.Size()

		snap.Files = append(snap.Files, fileName)
		snap.PackNumbers[fileName] = []int{}
		snap.Offsets[fileName] = []int64{}
		snap.Lengths[fileName] = []int64{}

		for fileBytesRemaining > 0 {
			copyBytes := min64(packBytesRemaining, fileBytesRemaining)
			bytesCopied, err := io.CopyN(packFile, in, copyBytes)

			if err == nil {
				fileBytesRemaining -= bytesCopied
				packBytesRemaining -= bytesCopied
				packOffset += bytesCopied

				snap.PackNumbers[fileName] = append(snap.PackNumbers[fileName], packCount)
				snap.Offsets[fileName] = append(snap.Offsets[fileName], packOffset)
				snap.Lengths[fileName] = append(snap.Lengths[fileName], bytesCopied)
			} else {
				if VerboseMode {
					fmt.Printf("Error writing file %s to pack %s, aborting\n", fileName, packPath)
				}

				return err
			}

			if packBytesRemaining <= 0 {
				packFile.Close()
				packCount++
				packOffset = 0
				packBytesRemaining = maxPackBytes
				packPath = filepath.Join(archiveFolder, fmt.Sprintf("pack%d.dat", packCount))
				var err error
				packFile, err = os.Create(packPath)

				if err != nil {
					if VerboseMode {
						fmt.Printf("Error creating pack file %s\n", packPath)
					}

					return err
				}

				if VerboseMode {
					fmt.Printf("Creating new pack file %s\n", packPath)
				}
			}
		}

		return nil
	}

	// fmt.Printf("No changes detected in %s for commit %s\n", workDir, snapshot.ID)
	filepath.Walk(workingDirectory, VersionFile)
	packFile.Close()
	snap.Write(archiveFolder)
	return nil
}

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

// func ListArchiveContents(archivePrefix string) {
// 	snap, err := ReadArchive(archivePrefix)

// 	if err != nil {
// 		fmt.Println("Error reading archive contents")
// 	}

// 	for i, file := range snap.Files {
// 		mtime := snap.ModTimes[file]
// 		bytes := snap.Lengths[file]

// 		if VerboseMode {
// 			fmt.Printf("%03d: %19s, %4d MB, %s\n", i, mtime, bytes/1000000, file)
// 		} else {
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

func AddOptionFlags(fs *flag.FlagSet) {
	fs.BoolVar(&VerboseMode, "verbose", false, "verbose mode")
	fs.BoolVar(&VerboseMode, "v", false, "verbose mode")
}

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
		AddOptionFlags(storeCmd)
		storeCmd.Parse(os.Args[2:])
		archiveFolder := storeCmd.Arg(0)
		inputFolder := storeCmd.Arg(1)
		StoreFolder(archiveFolder, inputFolder, 10*1024*1024)
		// } else if cmd == "list" || cmd == "ls" || cmd == "l" {
		// 	AddOptionFlags(listCmd)
		// 	listCmd.Parse(os.Args[2:])
		// 	archivePrefix := listCmd.Arg(0)
		// 	ListArchiveContents(archivePrefix)
		// } else if cmd == "extract" || cmd == "ex" || cmd == "e" || cmd == "x" {
		// 	AddOptionFlags(extractCmd)
		// 	extractCmd.Parse(os.Args[2:])
		// 	archivePrefix := extractCmd.Arg(0)
		// 	outputFolder := extractCmd.Arg(1)
		// 	ExtractArchive(archivePrefix, outputFolder)
	} else {
		fmt.Println("Unknown subcommand")
		os.Exit(1)
	}
}
