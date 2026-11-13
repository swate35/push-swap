package main

import (
	"math/rand"
	"testing"
)

// run applies the instructions on a copy of numbers and returns the stacks.
func run(t *testing.T, numbers []int, instructions []string) *stacks {
	t.Helper()
	s := &stacks{a: append([]int{}, numbers...), b: []int{}}
	for _, in := range instructions {
		switch in {
		case "sa":
			swap(s.a)
		case "sb":
			swap(s.b)
		case "ss":
			swap(s.a)
			swap(s.b)
		case "pa":
			s.b, s.a = push(s.b, s.a)
		case "pb":
			s.a, s.b = push(s.a, s.b)
		case "ra":
			rotate(s.a)
		case "rb":
			rotate(s.b)
		case "rr":
			rotate(s.a)
			rotate(s.b)
		case "rra":
			reverseRotate(s.a)
		case "rrb":
			reverseRotate(s.b)
		case "rrr":
			reverseRotate(s.a)
			reverseRotate(s.b)
		default:
			t.Fatalf("unknown instruction %q", in)
		}
	}
	return s
}

func TestSolveSorts(t *testing.T) {
	for size := 1; size <= 200; size++ {
		for try := 0; try < 20; try++ {
			numbers := rand.Perm(size * 3)[:size]
			for i := range numbers {
				numbers[i] -= size // include negatives
			}
			s := run(t, numbers, solve(numbers))
			if !sortstack(s.a) || len(s.b) != 0 {
				t.Fatalf("not sorted for %v: a=%v b=%v", numbers, s.a, s.b)
			}
		}
	}
}

func TestSolveCount(t *testing.T) {
	for size, max := range map[int]int{3: 2, 5: 12, 100: 1500, 500: 11500} {
		numbers := rand.Perm(size)
		if got := len(solve(numbers)); got > max {
			t.Errorf("%d numbers: %d instructions (max %d)", size, got, max)
		}
	}
}

func TestParseArgs(t *testing.T) {
	if _, err := parseArgs([]string{"1 one 3"}); err == nil {
		t.Error("expected error for non integer")
	}
	if _, err := parseArgs([]string{"1 2 1"}); err == nil {
		t.Error("expected error for duplicate")
	}
	got, err := parseArgs([]string{"3 -1", "2"})
	if err != nil || len(got) != 3 || got[0] != 3 || got[1] != -1 || got[2] != 2 {
		t.Errorf("unexpected result %v, %v", got, err)
	}
}
