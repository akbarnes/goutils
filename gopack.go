package main

import (
	"bufio"
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

func ExtractArchive(archiveFolder string, outputFolder string) error {
	snap, err := ReadSnapshot(archiveFolder)

	if err != nil {
		fmt.Println("Error reading archive contents")
	}

	if err := os.MkdirAll(outputFolder, 0777); err != nil {
		if VerboseMode {
			fmt.Printf("Error creating output folder %s\n", outputFolder)
		}

		return err
	}

	for _, file := range snap.Files {
		outPath := filepath.Join(outputFolder, file)

		if err := ExtractFile(file, archiveFolder, snap, outPath); err == nil {
			fmt.Println(file)
		}
		// else {
		// 	if VerboseMode {
		// 		fmt.Printf("Error extracting %s to %s\n", file, outPath)
		// 	}
		// }
	}

	return nil
}

func ExtractFile(file string, archiveFolder string, snap Snapshot, dst string) error {
	outputFolder := filepath.Dir(dst)

	if err := os.MkdirAll(outputFolder, 0777); err != nil {
		if VerboseMode {
			fmt.Printf("Error creating output folder %s\n", outputFolder)
		}

		return err
	}

	out, err := os.Create(dst)

	if err != nil {
		if VerboseMode {
			fmt.Printf("Error creating destination file %s:\n", dst)
			fmt.Println(err)
			fmt.Printf("\n")
		}

		return err
	}

	defer out.Close()
	w := bufio.NewWriter(out)

	for i, packNum := range snap.PackNumbers[file] {
		offset := snap.Offsets[file][i]
		nbytes := snap.Lengths[file][i]
		packPath := filepath.Join(archiveFolder, fmt.Sprintf("pack%d.dat", packNum))
		in, err := os.Open(packPath)
		defer in.Close()

		if err != nil {
			if VerboseMode {
				fmt.Println("Could not open pack file %s", packPath)
			}

			return err
		}

		if _, err := in.Seek(offset, 0); err != nil {
			if VerboseMode {
				fmt.Printf("Error seeking on pack file %s\n", packPath)
			}

			return err
		}

		// if _, err := io.CopyN(out, in, nbytes); err != nil {

		if wb, err := w.WriteString("buffered\n"); err != nil {
			if VerboseMode {
				fmt.Printf("\nError copying %d bytes from pack %s starting at %d bytes to destination file %s\n", nbytes, packPath, offset, dst)
				fmt.Println(err)
				fmt.Printf("\n")
			}

			return err
		} else {
			if VerboseMode {
				fmt.Printf("Wrote %d bytes\n", wb)
			}
		}

	}

	w.Flush()
	return nil
}

func Sum64(a []int64) int64 {
	var s int64 = 0

	for _, x := range a {
		s += x
	}

	return s
}

func ListArchiveContents(archiveFolder string) error {
	snap, err := ReadSnapshot(archiveFolder)

	if err != nil {
		fmt.Println("Error reading archive contents")
		return err
	}

	for i, file := range snap.Files {
		packs := len(snap.PackNumbers[file])
		MB := float64(Sum64(snap.Lengths[file])) / 1000000.0

		if VerboseMode {
			fmt.Printf("%3d: %3d chunks, %5.1f MB, %s\n", i+1, packs, MB, file)
		} else {
			fmt.Println(file)
		}
	}

	return nil
}

// Read a snapshot given a file path
func ReadSnapshot(archiveFolder string) (Snapshot, error) {
	snapshotPath := filepath.Join(archiveFolder, "snapshot.json")

	var mySnapshot Snapshot
	f, err := os.Open(snapshotPath)

	if err != nil {
		// panic(fmt.Sprintf("Error: Could not read snapshot file %s", snapshotPath))
		return Snapshot{}, err
	}

	defer f.Close()
	myDecoder := json.NewDecoder(f)

	if err := myDecoder.Decode(&mySnapshot); err != nil {
		fmt.Printf("Error:could not decode snapshot file %s\n", snapshotPath)
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
	extractCmd := flag.NewFlagSet("extract", flag.ExitOnError)

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
	} else if cmd == "list" || cmd == "ls" || cmd == "l" {
		AddOptionFlags(listCmd)
		listCmd.Parse(os.Args[2:])
		archiveFolder := listCmd.Arg(0)
		ListArchiveContents(archiveFolder)
	} else if cmd == "extract" || cmd == "ex" || cmd == "e" || cmd == "x" {
		AddOptionFlags(extractCmd)
		extractCmd.Parse(os.Args[2:])
		archiveFolder := extractCmd.Arg(0)
		outputFolder := extractCmd.Arg(1)
		ExtractArchive(archiveFolder, outputFolder)
	} else {
		fmt.Println("Unknown subcommand")
		os.Exit(1)
	}
}
