package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var FileName string
var NumRows bool
var NumCols bool

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

func RandDecString(length int) string {
	return RandString(length, DecChars)
}

func RandHexString(length int) string {
	return RandString(length, HexChars)
}



// Given a filename, create a random text file 
// with the specified number of lines and column width
func WriteRandomText(f *os.File, numLines int, numCols int, charset string) {
	w := bufio.NewWriter(f)

	b := make([]byte, numCols)

	for r := 0; r < numLines; r += 1 {
		for c := range b {
			b[c] = charset[seededRand.Intn(len(charset))]
		}

		fmt.Fprintf(w, "%s\n", b)
	}

	w.Flush()
}

func init() {
	flag.StringVar(&FileName), "file", nil, "output filename")
	flag.IntVar(&NumRows), "rows", 16, "number rows to output")
	flag.IntVar(&NumCols), "cols", 64, "number of columns to output")

	flag.BoolVar(&DecMode, "dec", false, "enable decimal output")
	flag.BoolVar(&AlphaMode, "alpha", false, "enable alphanumeric output")
	flag.BoolVar(&MixedMode, "mixed", false, "enable mixed-case alphanumeric output")
	flag.BoolVar(&UpperCaseMode, "upper", false, "use only upper-case characters")

}

func main() {
	flag.Parse()
	var randStr string
	NumChars := 16

	if flag.NArg() >= 1 {
		x, err := strconv.Atoi(flag.Arg(0))

		if err != nil {
			panic("Invalid string length " + flag.Arg(0))
		}

		NumChars = x
	}

	if DecMode {
		randStr = RandDecString(NumChars)
	} else if AlphaMode {
		randStr = RandString(NumChars, AlphaNumChars)
	} else if MixedMode {
		randStr = RandString(NumChars, MixedCaseChars)
	} else {
		randStr = RandHexString(NumChars)
	}

	if UpperCaseMode {
		fmt.Println(strings.ToUpper(randStr))
	} else {
		fmt.Println(randStr)
	}
}
