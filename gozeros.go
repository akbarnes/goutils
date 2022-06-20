package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
)

// Given a filename, create a random text file
// with the specified number of lines and column width
func WriteZeroBytes(MB int, fileName string) {
	fo, err := os.Create(fileName)
	defer fo.Close()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to open %s for writing, aborting\n", fileName)
		os.Exit(1)
	}

	bytes := make([]uint8, 1000*1000)

	for i := 0; i < MB; i++ {
		if err := binary.Write(fo, binary.LittleEndian, bytes); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write bytes, aborting: %v\n", err)
			os.Exit(1)
		}

		if i%1000 == 0 {
			fmt.Fprintf(os.Stderr, "%06d GB\n", i)
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
		fmt.Fprintf(os.Stderr, "Could not parse number of MB: %v\n", err)
		os.Exit(1)
	}

	FileName := os.Args[2]

	WriteZeroBytes(MB, FileName)
}
