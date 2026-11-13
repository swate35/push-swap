package main

import (
	"errors"
	"strconv"
	"strings"
)

// parseArgs converts the arguments into integers.
// It returns an error if a value is not an integer or appears twice.
func parseArgs(args []string) ([]int, error) {
	result := []int{}
	seen := map[int]bool{}

	for _, arg := range args {
		for _, field := range strings.Fields(arg) {
			n, err := strconv.Atoi(field)
			if err != nil {
				return nil, err
			}
			if seen[n] {
				return nil, errors.New("duplicate value")
			}
			seen[n] = true
			result = append(result, n)
		}
	}
	return result, nil
}

// sortstack reports whether the stack is sorted in ascending order.
func sortstack(a []int) bool {
	for i := 0; i < len(a)-1; i++ {
		if a[i] > a[i+1] {
			return false
		}
	}
	return true
}
