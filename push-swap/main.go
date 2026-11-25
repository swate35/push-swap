package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		return
	}

	numbers, err := parseArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error")
		return
	}

	for _, instruction := range solve(numbers) {
		fmt.Println(instruction)
	}
}
