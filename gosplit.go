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

func JoinFile(outputPath string, inputPrefix string) error {
	files, err := filepath.Glob(inputPrefix + ".*.part")

	if err != nil {
		if VerboseMode {
			fmt.Printf("Error listing parts for %s\n", inputPrefix)
		}

		return err
	}

	out, err := os.Create(outputPath)

	if err != nil {
		return err
	}

	for i := 1; i <= len(files); i++ {
		src := fmt.Sprintf("%s.%d.part", inputPrefix, i)
		fmt.Println(src)

		in, err := os.Open(src)

		if err != nil {
			if VerboseMode {
				fmt.Printf("Error opening part file %s:\n", src)
				fmt.Println(err)
			}

			return err
		}

		defer in.Close()
		io.Copy(out, in)

	}

	return out.Close()
}

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

		JoinFile(outputPath, inputPrefix)
	} else {
		fmt.Println("Unknown subcommand")
		os.Exit(1)
	}
}
