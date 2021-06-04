package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var FileName string
var NumRows int
var NumCols int

var DecMode bool
var AlphaMode bool
var MixedMode bool
var UpperCaseMode bool

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const DecChars = "0123456789"
const HexChars = DecChars + "abcdef"
const AlphaNumChars = DecChars + "abcdefghijklmnopqrstuvwxyz"
const MixedCaseChars = AlphaNumChars + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const UpperCaseHexChars = DecChars + "ABCDEF"
const UpperCaseAlphaNumChars = AlphaNumChars + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// Given a filename, create a random text file 
// with the specified number of lines and column width
func WriteRandomText(numLines int, numCols int, charset string) {
	b := make([]byte, numCols)

	for r := 0; r < numLines; r += 1 {
		for c := range b {
			b[c] = charset[seededRand.Intn(len(charset))]
		}

		fmt.Printf("%s\n", b)
	}
}

func init() {
	flag.IntVar(&NumRows, "rows", 16, "number rows to output")
	flag.IntVar(&NumCols, "cols", 64, "number of columns to output")

	flag.BoolVar(&DecMode, "dec", false, "enable decimal output")
	flag.BoolVar(&AlphaMode, "alpha", false, "enable alphanumeric output")
	flag.BoolVar(&MixedMode, "mixed", false, "enable mixed-case alphanumeric output")
	flag.BoolVar(&UpperCaseMode, "upper", false, "use only upper-case characters")

}

func main() {
	flag.Parse()

	if DecMode {
		WriteRandomText(NumRows, NumCols, DecChars)
	} else if AlphaMode {
		if UpperCaseMode {
			WriteRandomText(NumRows, NumCols, UpperCaseAlphaNumChars)
		} else {
			WriteRandomText(NumRows, NumCols, AlphaNumChars)
		}
	} else if MixedMode {
		WriteRandomText(NumRows, NumCols, MixedCaseChars)
	} else {
		if UpperCaseMode {
			WriteRandomText(NumRows, NumCols, UpperCaseHexChars)
		} else {
			WriteRandomText(NumRows, NumCols, HexChars)
		}
	} 
}

