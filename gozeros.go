package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

// Given a filename, create a random text file
// with the specified number of lines and column width
func WriteZeroBytes(numBytes int, fileName string) {
	fo, err := os.Create(fileName)
	defer fo.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s for writing, aborting\n", fileName)
		os.Exit(1)
	}

	for i := 0; i < numBytes; i++ {
		if err := binary.Write(fo, binary.LittleEndian, uint8(0)); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write byte, aborting: %v\n", err)
			os.Exit(1)
		}

		if i%(1000*1000) == 0 {
			MB := i / (1000 * 1000)
			fmt.Fprintf(os.Stderr, "%03d MB\n", MB)
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: gozeros <MB> <filename>\n")
		os.Exit(1)
	}

	MB, err := strconv.Atoi(os.Args[1])

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not parse number of bytes: %v\n", err)
		os.Exit(1)
	}

	NumBytes := 1000 * 1000 * MB
	FileName := os.Args[2]

	WriteZeroBytes(NumBytes, FileName)
}
