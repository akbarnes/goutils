package main

import (
	"os"
	"fmt"
	"io/ioutil"
)


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


func main() {
	diff, err := DiffFile(os.Args[1], os.Args[2])

	if err != nil {
		fmt.Println("Error comparing files")
	} else if diff {
		fmt.Println("Files are different")
	} else {
		fmt.Println("Files are the same")
	}
}