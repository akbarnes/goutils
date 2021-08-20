package main

import (
	"flag"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var Sacred string
var Variant string
var PrintJdupesCommand bool

// var RunJdupes bool
var OutputLanguage string

func init() {
	flag.StringVar(&Sacred, "sacred", "", "specify folder to preserve")
	flag.StringVar(&Variant, "variant", "", "specify folder to prune")
	flag.BoolVar(&PrintJdupesCommand, "command", false, "print jdupes command-line to run given sacred and variant folders")
	flag.StringVar(&OutputType, "language", "sh", "specify output language")
	// flag.BoolVar(&RunJdupes, "run", false, "run jdupes")
}

func main() {
	flag.Parse()
	// jdupesOutput := ""
	//
	// if flag.NArg() >= 1 {
	// 	jdupesOutput = flag.Arg(0)
	// }

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter text: ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

}
