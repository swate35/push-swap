package main

import (
	"math/rand"
	"testing"
)

func run(t *testing.T, numbers []int, instructions []string) *stacks {
	return runFrom(t, numbers, nil, instructions)
}

func runFrom(t *testing.T, aValues, bValues []int, instructions []string) *stacks {
	t.Helper()
	s := &stacks{a: append([]int{}, aValues...), b: append([]int{}, bValues...)}
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
				numbers[i] -= size
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
	if _, err := parseArgs([]string{""}); err == nil {
		t.Error("expected error for empty argument")
	}
	if _, err := parseArgs([]string{"   "}); err == nil {
		t.Error("expected error for whitespace-only argument")
	}
	if _, err := parseArgs([]string{"1 2 1"}); err == nil {
		t.Error("expected error for duplicate")
	}
	if _, err := parseArgs([]string{"2147483648"}); err == nil {
		t.Error("expected error for value above int32")
	}
	if _, err := parseArgs([]string{"-2147483649"}); err == nil {
		t.Error("expected error for value below int32")
	}
	if _, err := parseArgs([]string{"2147483647 -2147483648"}); err != nil {
		t.Errorf("expected int32 boundary values to be accepted, got %v", err)
	}
	got, err := parseArgs([]string{"3 -1", "2"})
	if err != nil || len(got) != 3 || got[0] != 3 || got[1] != -1 || got[2] != 2 {
		t.Errorf("unexpected result %v, %v", got, err)
	}
}

type referenceNode struct {
	a, b     []int
	distance int
}

func referenceStateKey(a, b []int) string {
	key := make([]byte, 0, len(a)+len(b)+1)
	for _, value := range a {
		key = append(key, byte(value))
	}
	key = append(key, 0xff)
	for _, value := range b {
		key = append(key, byte(value))
	}
	return string(key)
}

func referenceDistances(t *testing.T, size int) map[string]int {
	goal := make([]int, size)
	for i := range goal {
		goal[i] = i
	}

	queue := []referenceNode{{a: goal}}
	distances := map[string]int{referenceStateKey(goal, nil): 0}
	for head := 0; head < len(queue); head++ {
		current := queue[head]
		for _, instruction := range smallInstructions {
			next := runFrom(t, current.a, current.b, []string{instruction})
			key := referenceStateKey(next.a, next.b)
			if _, visited := distances[key]; visited {
				continue
			}
			distance := current.distance + 1
			distances[key] = distance
			queue = append(queue, referenceNode{a: next.a, b: next.b, distance: distance})
		}
	}
	return distances
}

func TestSolveIsOptimalForUpToFiveElements(t *testing.T) {
	for size := 1; size <= 5; size++ {
		distances := referenceDistances(t, size)
		numbers := make([]int, size)
		for i := range numbers {
			numbers[i] = i
		}

		var checkPermutations func(int)
		checkPermutations = func(position int) {
			if position == size {
				want := distances[referenceStateKey(numbers, nil)]
				if got := len(solve(numbers)); got != want {
					t.Errorf("solve(%v) uses %d instructions, shortest is %d", numbers, got, want)
				}
				return
			}

			for i := position; i < size; i++ {
				numbers[position], numbers[i] = numbers[i], numbers[position]
				checkPermutations(position + 1)
				numbers[position], numbers[i] = numbers[i], numbers[position]
			}
		}
		checkPermutations(0)
	}
}
