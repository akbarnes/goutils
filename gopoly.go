package main

import (
	"fmt"
    "github.com/restic/chunker"
)

func main() {
    poly, err := chunker.RandomPolynomial()

	if err != nil {
		panic(fmt.Sprintf("Error generating polynomial: %v\n", err))
	}

    fmt.Printf("%v\n", poly)
}
