package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const FeetToMeters = 0.3048
const InchesToMeters = 0.0254
const MilesToMeters = 1609.34

// examples
// 1     2     3   4
// 15.3 inches to feet
// 32.7 degrees to radians

func expandUnits(unit string) string {
	unit = strings.ToLower(unit)

	switch unit {
	case "ft":
		return "feet"
	case "m":
		return "meters"
	case "in":
		return "inches"
	case "deg":
		return "degrees"
	case "f":
		return "farenheight"
	case "c":
		return "celsius"
	case "centigrade":
		return "celsius"
	}

	return unit
}

func main() {
	x, err := strconv.ParseFloat(os.Args[1], 64)

	if err != nil {
		panic("Invalid value " + os.Args[1])
	}

	from := expandUnits(os.Args[2])
	to := expandUnits(os.Args[4])
	// fmt.Printf("Converting %f from %s to %s\n", x, from, to)

	switch from {
	case "feet":
		x *= FeetToMeters
	case "inches":
		x *= InchesToMeters
	case "miles":
		x *= MilesToMeters
	case "degrees":
		x *= (math.Pi / 180.0)
	case "farenheight":
		x = (x - 32.0) * 5 / 9
	}

	switch to {
	case "feet":
		x /= FeetToMeters
	case "inches":
		x /= InchesToMeters
	case "miles":
		x /= MilesToMeters
	case "degrees":
		x /= (math.Pi / 180.0)
	case "farenheight":
		x = (9/5)*x + 32
	}

	fmt.Printf("%0.6f\n", x)
}
