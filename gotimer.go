package main

import (
	"fmt"
	"os"
	// "strconv"
	"time"

	"github.com/sqweek/dialog"
)

func main() {
	var err error

	duration, err := time.ParseDuration(os.Args[1])
	if err != nil {
		panic("error: invalid duration: " +  os.Args[1])
	}

	secs := int(duration.Seconds())
	// fmt.Println(secs)

	if err != nil {
		fmt.Println(os.Args)
		panic("Invalid time specification")
	}

	for secs > 0 {
		mm := secs / 60
		ss := secs - 60*mm
		fmt.Printf("\r%02d:%02d", mm, ss)
		secs -= 1
		time.Sleep(time.Second)
	}

	fmt.Printf("\r%02d:%02d\n", 0, 0)
	dialog.Message("%s", "Time's Up!").Title("GoTime").Info()
}
