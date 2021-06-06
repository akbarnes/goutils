package main

import (
	// "os"
	"flag"
	"fmt"
	"io/ioutil"
)

var HumanReadable bool
var ShowErrors bool
var Json bool

func DiffFile(f1 string, f2 string) (bool, error) {
	dat1, err1 := ioutil.ReadFile(f1)

	if err1 != nil {
		return true, err1
	}

	dat2, err2 := ioutil.ReadFile(f2)

	if err2 != nil {
		return true, err2
	}	

	if len(dat1) != len(dat2) {
		return true, nil
	}

	for i := 0; i < len(dat1); i++ {
		x1 := dat1[i]
		x2 := dat2[i]

		if x1 != x2 {
			return true, nil
		}
	}

	return false, nil
}

func init() {
	flag.BoolVar(&HumanReadable, "human", false, "Show human-readable output")
	flag.BoolVar(&ShowErrors, "errors", false, "Report errors reading files")
	flag.BoolVar(&Json, "json", false, "Use json output")

}

func main() {
	flag.Parse()
	diff, err := DiffFile(flag.Arg(0), flag.Arg(1))

	if HumanReadable {
		if err != nil {
			fmt.Println("Error comparing files")
		} else if diff {
			fmt.Println("Files are different")
		} else {
			fmt.Println("Files are equal")
		}
	} else if Json {
		if err != nil {
			if ShowErrors {
				fmt.Println("null")
			} else {
				fmt.Println("true")
			}
		} else if diff {
			fmt.Println("true")
		} else {
			fmt.Println("false")
		}
	} else {
		if err != nil {
			if ShowErrors {
				fmt.Println("error")
			} else {
				fmt.Println("different")
			}
		} else if diff {
			fmt.Println("different")
		} else {
			fmt.Println("equal")
		}
	}
}