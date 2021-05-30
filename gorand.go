package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var GroupSize int
var NumChars int
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

func RandString(length int, charset string) string {
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(b)
}

func init() {
	flag.BoolVar(&DecMode, "dec", false, "enable decimal output")
	flag.BoolVar(&AlphaMode, "alpha", false, "enable alphanumeric output")
	flag.BoolVar(&MixedMode, "mixed", false, "enable mixed-case alphanumeric output")
	flag.BoolVar(&UpperCaseMode, "upper", false, "use only upper-case characters")
	flag.IntVar(&GroupSize, "group", 4, "group size when splitting with dashes")
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
		randStr = strings.ToUpper(randStr)
	}

	for i, c := range randStr {
		fmt.Printf("%c", c)

		if i < len(randStr) - 1 && GroupSize > 0 && (i + 1) % GroupSize == 0 {
			fmt.Printf("-")
		}
	}

	fmt.Println("")
}
