package main

import (
	"errors"
	"strconv"
	"strings"
)

func parseArgs(args []string) ([]int, error) {
	result := []int{}
	seen := map[int]bool{}

	for _, arg := range args {
		fields := strings.Fields(arg)
		if len(fields) == 0 {
			return nil, errors.New("empty argument")
		}

		for _, field := range fields {
			value, err := strconv.ParseInt(field, 10, 32)
			if err != nil {
				return nil, err
			}
			n := int(value)
			if seen[n] {
				return nil, errors.New("duplicate value")
			}
			seen[n] = true
			result = append(result, n)
		}
	}
	return result, nil
}

func sortstack(a []int) bool {
	for i := 0; i < len(a)-1; i++ {
		if a[i] > a[i+1] {
			return false
		}
	}
	return true
}
