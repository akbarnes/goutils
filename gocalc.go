package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func printStack(s []float64) {
	fmt.Printf("T: %0.6f\n", s[3])
	fmt.Printf("Z: %0.6f\n", s[2])
	fmt.Printf("Y: %0.6f\n", s[1])
	fmt.Printf("a: %0.6f\n", s[0])
}

func pushStack(s []float64, a float64) {
	s[3] = s[2]
	s[2] = s[1]
	s[1] = s[0]
	s[0] = a
}

func dupStack(s []float64) {
	pushStack(s, s[0])
}

func rotateStack(s []float64) {
	s[0] = s[1]
	s[1] = s[2]
	s[2] = s[3]
}

func swapStack(s []float64) {
	a := s[0]
	s[0] = s[1]
	s[1] = a
}

func clearStack(s []float64) {
	for i := 0; i <= 3; i++ {
		s[i] = 0.0
	}
}

func rotateBottom(s []float64) {
	s[1] = s[2]
	s[2] = s[3]
}

func addStack(s []float64) {
	s[0] += s[1]
	rotateBottom(s)
}

func multiplyStack(s []float64) {
	s[0] *= s[1]
	rotateBottom(s)
}

func subtractStack(s []float64) {
	s[0] = s[1] - s[0]
	rotateBottom(s)
}

func divideStack(s []float64) {
	s[0] = s[1] / s[0]
	rotateBottom(s)
}

func powerStack(s []float64) {
	s[0] = math.Pow(s[1], s[0])
	rotateBottom(s)
}

// HMS
func main() {
	var s = []float64{0, 0, 0, 0}
	var input = ""

	for input != "quit" && input != "q" {
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		input, _ = reader.ReadString('\n')
		// fmt.Printf("<%s>\n", input)
		input = input[0 : len(input)-2]
		// fmt.Printf("<%s>\n", input)
		fmt.Println("")

		if input == "print" || input == "p" {
			printStack(s)
		} else if input == "rotate" || input == "r" {
			rotateStack(s)
			printStack(s)
		} else if input == "swap" || input == "s" || input == "w" {
			swapStack(s)
			printStack(s)
		} else if input == "+" || input == "addStack" {
			addStack(s)
			printStack(s)
		} else if input == "*" || input == "multiplyStack" || input == "mult" {
			multiplyStack(s)
			printStack(s)
		} else if input == "-" || input == "subtractStack" || input == "sub" {
			subtractStack(s)
			printStack(s)
		} else if input == "/" || input == "divideStack" || input == "div" {
			divideStack(s)
			printStack(s)
		} else if input == "sine" || input == "sin" || input == "s" {
			s[0] = math.Sin(s[0])
			printStack(s)
		} else if input == "arcsine" || input == "asin" || input == "S" {
			s[0] = math.Asin(s[0])
			printStack(s)
		} else if input == "cosine" || input == "cos" || input == "c" {
			s[0] = math.Cos(s[0])
			printStack(s)
		} else if input == "arccosine" || input == "acos" || input == "C" {
			s[0] = math.Acos(s[0])
			printStack(s)
		} else if input == "tangent" || input == "tan" || input == "t" {
			s[0] = math.Tan(s[0])
			printStack(s)
		} else if input == "arctangent" || input == "atan" || input == "T" {
			s[0] = math.Atan(s[0])
			printStack(s)
		} else if input == "invert" || input == "inv" || input == "i" {
			s[0] = 1.0 / s[0]
			printStack(s)
		} else if input == "square_root" || input == "sqrt" {
			s[0] = math.Sqrt(s[0])
			printStack(s)
		} else if input == "square" || input == "sqr" {
			s[0] = s[0] * s[0]
			printStack(s)
		} else if input == "ln" || input == "l" {
			s[0] = math.Log(s[0])
			printStack(s)
		} else if input == "log" || input == "L" {
			s[0] = math.Log(s[0]) / math.Log(10.0)
			printStack(s)
		} else if input == "exp" || input == "e" {
			s[0] = math.Exp(s[0])
			printStack(s)
		} else if input == "+/-" || input == "negate" || input == "n" {
			s[0] = -s[0]
			printStack(s)
		} else if input == "pi" || input == "p" {
			pushStack(s, math.Pi)
			printStack(s)
		} else if input == "y^x" || input == "y" {
			powerStack(s)
			printStack(s)
		} else if input == "10^x" {
			s[0] = math.Pow(10, s[0])
			printStack(s)
		} else if input == "clear" || input == "clr" {
			clearStack(s)
			printStack(s)
		} else if input == "duplicate" || input == "dup" || input == "d" {
			dupStack(s)
			printStack(s)
		} else {
			a, err := strconv.ParseFloat(input, 64)

			if err != nil {
				fmt.Printf("Invalid input %s\n", input)
			}

			pushStack(s, a)
			printStack(s)
		}

	}

}
