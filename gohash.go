package main

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Optional timestamp
	t := time.Now()

	fmt.Printf("File: %s\n", os.Args[1])
	fmt.Println(t.Format("2006-01-02T15-04-05"))

	data, err := ioutil.ReadFile(os.Args[1])
	check(err)
	sum := sha256.Sum256(data)
	fmt.Printf("%x", sum)
}
