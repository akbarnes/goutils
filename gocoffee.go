package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var VerboseMode bool

func init() {
	flag.BoolVar(&VerboseMode, "v", false, "toggle verbose output")
}

func main() {
	flag.Parse()

	cups, err := strconv.ParseFloat(os.Args[1], 64)

	if err != nil {
		panic("Invalid number of cups specified %s " + os.Args[1])
	}

	if VerboseMode {
		fmt.Printf("%0.1f tbsp\n", 2*cups)
	} else {
		fmt.Println(2 * cups)
	}
}
