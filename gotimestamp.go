package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	t := time.Now()

	if len(os.Args) >= 2 {
		fmt.Println(t.Format("2006-01-02T15-04-05"))
	} else {
		fmt.Println(t.Format("2006-01-02"))
	}
}
