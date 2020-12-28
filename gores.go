package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func main() {
	ResistMap := map[string]float64{"black": 0, "brown": 1, "red": 2, "orange": 3, "yellow": 4, "green": 5, "blue": 6, "violet": 7, "gray": 8, "grey": 8, "white": 9}
	InvResist := []string{"black", "brown", "red", "orange", "yellow", "green", "blue", "violet", "gray", "white"}
	TolMap := map[string]int{"gold": 5, "silver": 10}
	UnitsMap := map[string]float64{"k": 1e3, "M": 1e6}

	var bands [4]string

	if len(os.Args) <= 3 {
		r, err := strconv.ParseFloat(os.Args[1], 64)

		if err != nil {
			panic("Invalid resistance " + os.Args[1])
		}

		if len(os.Args) == 3 {
			r *= UnitsMap[os.Args[2]]
		}

		exponent := 0

		if r >= 100.0 {
			for r >= 100.0 {
				r /= 10.0
				exponent += 1
			}
		}

		if r < 10 {
			fmt.Printf("black %s ", InvResist[int(r)])
		} else {
			d1 := r / 10.0
			d2 := r - 10.0*d1

			fmt.Printf("%s %s ", InvResist[int(d1)], InvResist[int(d2)])
		}

		fmt.Printf("%s\n", InvResist[exponent])
	} else {
		bands[0] = os.Args[1]
		bands[1] = os.Args[2]
		bands[2] = os.Args[3]
		bands[3] = ""

		if len(os.Args) >= 5 {
			bands[3] = os.Args[4]
		}

		r := 10*ResistMap[bands[0]] + ResistMap[bands[1]]
		r *= math.Pow(10, ResistMap[bands[2]])

		if r >= 1e6 {
			fmt.Printf("%0.1fM", r/1e6)
		} else if r >= 1e3 {
			fmt.Printf("%0.1fk", r/1e3)
		} else {
			fmt.Printf("%0.1f", r)
		}

		if len(bands[3]) > 0 {
			fmt.Printf(" +/- %d%%", TolMap[bands[3]])
		}

		fmt.Println("")
	}
}
