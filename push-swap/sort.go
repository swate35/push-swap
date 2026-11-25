package main

import (
	"math/bits"
	"sort"
)

func solve(numbers []int) []string {
	s := &stacks{a: append([]int{}, numbers...), b: []int{}}

	if !sortstack(s.a) {
		switch len(s.a) {
		case 2:
			s.sa()
		case 3:
			sort3(s)
		case 4, 5:
			sortSmall(s)
		default:
			radix(s)
		}
	}
	return optimize(s.instructions)
}

func sort3(s *stacks) {
	x, y, z := s.a[0], s.a[1], s.a[2]

	switch {
	case x > y && y < z && x < z:
		s.sa()
	case x > y && y > z:
		s.sa()
		s.rra()
	case x > y && y < z && x > z:
		s.ra()
	case x < y && y > z && x < z:
		s.sa()
		s.ra()
	case x < y && y > z && x > z:
		s.rra()
	}
}

func sortSmall(s *stacks) {
	for len(s.a) > 3 {
		pos := minIndex(s.a)
		if pos <= len(s.a)/2 {
			for i := 0; i < pos; i++ {
				s.ra()
			}
		} else {
			for i := 0; i < len(s.a)-pos; i++ {
				s.rra()
			}
		}
		s.pb()
	}

	if !sortstack(s.a) {
		sort3(s)
	}
	for len(s.b) > 0 {
		s.pa()
	}
}

func minIndex(a []int) int {
	min := 0
	for i, v := range a {
		if v < a[min] {
			min = i
		}
	}
	return min
}

func normalize(a []int) {
	sorted := append([]int{}, a...)
	sort.Ints(sorted)
	for i, v := range a {
		a[i] = sort.SearchInts(sorted, v)
	}
}

func radix(s *stacks) {
	normalize(s.a)
	nbBits := bits.Len(uint(len(s.a) - 1))

	for bit := 0; bit < nbBits; bit++ {
		size := len(s.a)
		for i := 0; i < size; i++ {
			if (s.a[0]>>bit)&1 == 0 {
				s.pb()
			} else {
				s.ra()
			}
		}
		for len(s.b) > 0 {
			s.pa()
		}
	}
}

var opposite = map[string]string{
	"pa": "pb", "pb": "pa",
	"ra": "rra", "rra": "ra",
	"rb": "rrb", "rrb": "rb",
}

var merged = map[string]string{
	"ra": "rb", "rb": "ra",
	"rra": "rrb", "rrb": "rra",
	"sa": "sb", "sb": "sa",
}

var mergedName = map[string]string{
	"ra": "rr", "rb": "rr",
	"rra": "rrr", "rrb": "rrr",
	"sa": "ss", "sb": "ss",
}

func optimize(instructions []string) []string {
	result := []string{}

	for _, instr := range instructions {
		if len(result) > 0 {
			last := result[len(result)-1]
			if opposite[last] == instr {
				result = result[:len(result)-1]
				continue
			}
			if merged[last] == instr {
				result[len(result)-1] = mergedName[last]
				continue
			}
		}
		result = append(result, instr)
	}
	return result
}
