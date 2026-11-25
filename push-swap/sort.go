package main

import (
	"math/bits"
	"sort"
)

const exactSortLimit = 7

var smallInstructions = [...]string{
	"pa", "pb", "sa", "sb", "ss", "ra", "rb", "rr", "rra", "rrb", "rrr",
}

type smallState struct {
	a, b       [exactSortLimit]uint8
	aLen, bLen uint8
}

type searchNode struct {
	state       smallState
	parent      int
	instruction uint8
}

func solve(numbers []int) []string {
	s := &stacks{a: append([]int{}, numbers...), b: []int{}}

	if sortstack(s.a) {
		return nil
	}
	if len(s.a) <= exactSortLimit {
		return exactSort(s.a)
	}

	radix(s)
	return optimize(s.instructions)
}

func exactSort(numbers []int) []string {
	sorted := append([]int{}, numbers...)
	sort.Ints(sorted)

	start := smallState{aLen: uint8(len(numbers))}
	goal := smallState{aLen: uint8(len(numbers))}
	for i, value := range numbers {
		start.a[i] = uint8(sort.SearchInts(sorted, value))
		goal.a[i] = uint8(i)
	}

	nodes := []searchNode{{state: start, parent: -1}}
	visited := map[smallState]int{start: 0}

	for head := 0; head < len(nodes); head++ {
		current := nodes[head].state
		for instruction := range smallInstructions {
			next, ok := applySmallInstruction(current, uint8(instruction))
			if !ok || next == current {
				continue
			}

			if _, exists := visited[next]; exists {
				continue
			}

			nextIndex := len(nodes)
			nodes = append(nodes, searchNode{
				state:       next,
				parent:      head,
				instruction: uint8(instruction),
			})
			visited[next] = nextIndex

			if next == goal {
				return reconstructInstructions(nodes, nextIndex)
			}
		}
	}

	panic("valid stack state has no sorting instruction sequence")
}

func applySmallInstruction(state smallState, instruction uint8) (smallState, bool) {
	next := state
	aLen, bLen := int(next.aLen), int(next.bLen)

	switch instruction {
	case 0: // pa
		if bLen == 0 {
			return next, false
		}
		copy(next.a[1:aLen+1], next.a[:aLen])
		next.a[0] = next.b[0]
		copy(next.b[:bLen-1], next.b[1:bLen])
		next.b[bLen-1] = 0
		next.aLen++
		next.bLen--
	case 1: // pb
		if aLen == 0 {
			return next, false
		}
		copy(next.b[1:bLen+1], next.b[:bLen])
		next.b[0] = next.a[0]
		copy(next.a[:aLen-1], next.a[1:aLen])
		next.a[aLen-1] = 0
		next.aLen--
		next.bLen++
	case 2: // sa
		swapSmall(&next.a, aLen)
	case 3: // sb
		swapSmall(&next.b, bLen)
	case 4: // ss
		swapSmall(&next.a, aLen)
		swapSmall(&next.b, bLen)
	case 5: // ra
		rotateSmall(&next.a, aLen)
	case 6: // rb
		rotateSmall(&next.b, bLen)
	case 7: // rr
		rotateSmall(&next.a, aLen)
		rotateSmall(&next.b, bLen)
	case 8: // rra
		reverseRotateSmall(&next.a, aLen)
	case 9: // rrb
		reverseRotateSmall(&next.b, bLen)
	case 10: // rrr
		reverseRotateSmall(&next.a, aLen)
		reverseRotateSmall(&next.b, bLen)
	default:
		return next, false
	}

	return next, true
}

func swapSmall(stack *[exactSortLimit]uint8, length int) {
	if length >= 2 {
		stack[0], stack[1] = stack[1], stack[0]
	}
}

func rotateSmall(stack *[exactSortLimit]uint8, length int) {
	if length < 2 {
		return
	}
	first := stack[0]
	copy(stack[:length-1], stack[1:length])
	stack[length-1] = first
}

func reverseRotateSmall(stack *[exactSortLimit]uint8, length int) {
	if length < 2 {
		return
	}
	last := stack[length-1]
	copy(stack[1:length], stack[:length-1])
	stack[0] = last
}

func reconstructInstructions(nodes []searchNode, index int) []string {
	instructions := make([]string, 0)
	for nodes[index].parent >= 0 {
		instructions = append(instructions, smallInstructions[nodes[index].instruction])
		index = nodes[index].parent
	}

	for left, right := 0, len(instructions)-1; left < right; left, right = left+1, right-1 {
		instructions[left], instructions[right] = instructions[right], instructions[left]
	}
	return instructions
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
